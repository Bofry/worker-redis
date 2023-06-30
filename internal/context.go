package internal

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Bofry/trace"
)

var (
	_ context.Context    = new(Context)
	_ trace.ValueContext = new(Context)
)

type Context struct {
	ConsumerGroup string
	ConsumerName  string

	context context.Context
	logger  *log.Logger

	invalidMessageHandler MessageHandler

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
func (*Context) Err() error {
	return nil
}

// Value implements context.Context.
func (c *Context) Value(key any) any {
	if key == nil {
		return nil
	}
	if c.values != nil {
		return c.values[key]
	}
	if c.context != nil {
		return c.context.Value(key)
	}
	return nil
}

// SetValue implements trace.ValueContext.
func (c *Context) SetValue(key interface{}, value interface{}) {
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

func (c *Context) ThrowInvalidMessageError(message *Message) {
	if c.invalidMessageHandler != nil {
		c.invalidMessageHandler.ProcessMessage(c, message)
	}
}

func (c *Context) clone() *Context {
	return &Context{
		ConsumerGroup:         c.ConsumerGroup,
		ConsumerName:          c.ConsumerName,
		context:               c.context,
		logger:                c.logger,
		invalidMessageHandler: c.invalidMessageHandler,
		values:                c.values,
	}
}
