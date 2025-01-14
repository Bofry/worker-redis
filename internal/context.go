package internal

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/trace"
)

var (
	_ context.Context    = new(Context)
	_ trace.ValueContext = new(Context)
)

type Context struct {
	ConsumerGroup string
	ConsumerName  string

	consumer *redis.Consumer

	context        context.Context // parent context
	err            error
	errExisted     int32
	logger         *log.Logger
	disableLogging bool
	aborted        bool

	invalidMessageHandler MessageHandler
	invalidMessageSent    int32

	values     map[interface{}]interface{}
	valuesOnce sync.Once
}

// Deadline implements context.Context.
func (*Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context.Context.
func (*Context) Done() <-chan struct{} {
	return nil
}

// Err implements context.Context.
func (c *Context) Err() error {
	if c.err != nil {
		return c.err
	}
	if c.context != nil {
		return c.context.Err()
	}
	return nil
}

func (c *Context) Break() {
	c.aborted = true
}

func (c *Context) IsAborted() bool {
	return c.aborted
}

func (c *Context) CatchErr(err error) {
	if err != nil {
		if c.err != nil {
			if atomic.CompareAndSwapInt32(&c.errExisted, 0, 1) {
				c.err = err
			}
		}
	}
}

// Value implements context.Context.
func (c *Context) Value(key any) any {
	if key == nil {
		return nil
	}
	if c.values != nil {
		v := c.values[key]
		if v != nil {
			return v
		}
	}
	if c.context != nil {
		return c.context.Value(key)
	}
	return nil
}

// SetValue implements trace.ValueContext.
func (c *Context) SetValue(key interface{}, value interface{}) {
	if c.aborted {
		return
	}

	if key == nil {
		return
	}
	if c.values == nil {
		c.valuesOnce.Do(func() {
			if c.values == nil {
				c.values = make(map[interface{}]interface{})
			}
		})
	}
	c.values[key] = value
}

func (c *Context) Logger() *log.Logger {
	return c.logger
}

func (c *Context) CanRecordingLog() bool {
	return !c.disableLogging
}

func (c *Context) RecordingLog(v bool) {
	if c.aborted {
		return
	}

	c.disableLogging = !v
}

func (c *Context) InvalidMessage(message *Message) {
	if c.aborted {
		return
	}

	if !atomic.CompareAndSwapInt32(&c.invalidMessageSent, 0, 1) {
		c.logger.Fatal("invalid operation; message has already been sent to InvalidMessageHandler")
	}

	GlobalContextHelper.InjectReplyCode(c, ABORT)

	if c.invalidMessageHandler != nil {
		var (
			tr       = GlobalTracerManager.GenerateManagedTracer(c.invalidMessageHandler)
			prevSpan = trace.SpanFromContext(c)
		)

		sp := tr.Start(prevSpan.Context(), __INVALID_MESSAGE_SPAN_NAME)
		defer func() {
			fmt.Println("(c *Context) InvalidMessage()")
			sp.End()
		}()

		ctx := &Context{
			logger:                c.logger,
			values:                c.values,
			context:               c,
			invalidMessageHandler: MessageHandleProc(StopRecursiveForwardMessageHandler),
		}
		trace.SpanToContext(ctx, sp)

		c.invalidMessageHandler.ProcessMessage(ctx, message)
	}
}

func (c *Context) Pause(streams ...string) error {
	return c.consumer.Pause(streams...)
}

func (c *Context) Resume(streams ...string) error {
	return c.consumer.Resume(streams...)
}

func (c *Context) Status() StatusCode {
	return GlobalContextHelper.ExtractReplyCode(c)
}

func (c *Context) clone() *Context {
	return &Context{
		ConsumerGroup:         c.ConsumerGroup,
		ConsumerName:          c.ConsumerName,
		consumer:              c.consumer,
		context:               c.context,
		logger:                c.logger,
		invalidMessageHandler: c.invalidMessageHandler,
		values:                c.values,
		aborted:               c.aborted,
	}
}
