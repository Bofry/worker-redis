package internal

import (
	"context"
	"fmt"

	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/lib-redis-stream/tracing"
	"github.com/Bofry/trace"
)

type MessageDispatcher struct {
	MessageHandleService   *MessageHandleService
	MessageTracerService   *MessageTracerService
	MessageObserverService *MessageObserverService
	Router                 Router

	OnHostErrorProc OnHostErrorHandler

	ErrorHandler          ErrorHandler
	InvalidMessageHandler MessageHandler

	StreamSet map[string]StreamOffset

	invalidMessageHandlerWrapper MessageHandler
}

func (d *MessageDispatcher) StreamOffsets() []StreamOffset {
	var (
		streams = d.StreamSet
	)

	offsets := make([]StreamOffset, 0, len(streams))
	for _, v := range streams {
		offsets = append(offsets, v)
	}
	return offsets
}

func (d *MessageDispatcher) Streams() []string {
	var (
		router = d.Router
	)

	if router != nil {
		keys := make([]string, 0, len(router))
		for k := range router {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *MessageDispatcher) ProcessMessage(ctx *Context, message *Message) {
	// start tracing
	var (
		handlerID = d.Router.FindHandlerComponentID(message.Stream)
		carrier   = tracing.NewMessageStateCarrier(&message.Content().State)

		spanName string = message.Stream
		tr       *trace.SeverityTracer
		sp       *trace.SeveritySpan
	)

	tr = d.MessageTracerService.Tracer(handlerID)
	sp = tr.ExtractWithPropagator(
		ctx,
		d.MessageTracerService.TextMapPropagator(),
		carrier,
		spanName)
	defer sp.End()

	processingState := ProcessingState{
		Stream: message.Stream,
		Tracer: tr,
		Span:   sp,
	}

	// set invalidMessageHandler
	ctx.invalidMessageHandler = d.invalidMessageHandlerWrapper

	// register observer into message
	d.MessageObserverService.RegisterMessageObservers(message, handlerID)

	d.MessageHandleService.ProcessMessage(ctx, message, processingState, new(Recover))
}

func (d *MessageDispatcher) subscribe(consumer *redis.Consumer) error {
	var (
		streams = d.Streams()
		offsets = make([]StreamOffsetInfo, 0, len(streams))
	)

	for _, stream := range streams {
		offsets = append(offsets, redis.Stream(stream))
	}

	return consumer.Subscribe(offsets...)
}

func (d *MessageDispatcher) internalProcessMessage(ctx *Context, message *Message, state ProcessingState, recover *Recover) {
	recover.
		Defer(func(err interface{}) {
			if err != nil {
				d.processError(ctx, message, err)
			}
		}).
		Do(func(finalizer Finalizer) {
			var (
				tr     *trace.SeverityTracer = state.Tracer
				sp     *trace.SeveritySpan   = state.Span
				stream string                = state.Stream
			)
			_ = tr

			// set Span
			trace.SpanToContext(ctx, sp)

			finalizer.Add(func(err interface{}) {
				if err != nil {
					if e, ok := err.(error); ok {
						sp.Err(e)
					} else if e, ok := err.(string); ok {
						sp.Err(fmt.Errorf(e))
					} else if e, ok := err.(fmt.Stringer); ok {
						sp.Err(fmt.Errorf(e.String()))
					} else {
						sp.Err(fmt.Errorf("%+v", err))
					}
				}

				var (
					reply = GlobalContextHelper.ExtractReplyCode(ctx)
				)

				switch reply {
				case PASS:
					sp.Reply(trace.PASS, reply)
				case FAIL, ABORT:
					sp.Reply(trace.FAIL, reply)
				}
			})

			sp.Tags(
				// TODO: add redis server version
				trace.Stream(stream),
				trace.ConsumerGroup(ctx.ConsumerGroup),
				trace.MessageID(message.ID),
			)

			handler := d.Router.Get(message.Stream)
			if handler != nil {
				handler.ProcessMessage(ctx, message)
				{
					reply := GlobalContextHelper.ExtractReplyCode(ctx)
					if reply == UNSET {
						GlobalContextHelper.InjectReplyCode(ctx, FAIL)
					}
				}
				return
			}
			ctx.InvalidMessage(message)
		})
}

func (d *MessageDispatcher) init() {
	// register the default MessageHandleModule
	stdMessageHandleModule := NewStdMessageHandleModule(d)
	d.MessageHandleService.Register(stdMessageHandleModule)

	if d.InvalidMessageHandler != nil {
		d.invalidMessageHandlerWrapper = &StdInvalidMessageHandler{
			invalidMessageHandler: d.InvalidMessageHandler,
		}
	}
}

func (d *MessageDispatcher) processError(ctx *Context, message *Message, err interface{}) {
	if d.ErrorHandler != nil {
		d.ErrorHandler(ctx, message, err)
	}
}

func (d *MessageDispatcher) start(ctx context.Context) {
	err := d.MessageHandleService.triggerStart(ctx)
	if err != nil {
		var disposed bool = false
		if d.OnHostErrorProc != nil {
			disposed = d.OnHostErrorProc(err)
		}
		if !disposed {
			RedisWorkerLogger.Fatalf("%+v", err)
		}
	}
}

func (d *MessageDispatcher) stop(ctx context.Context) {
	for err := range d.MessageHandleService.triggerStop(ctx) {
		if err != nil {
			RedisWorkerLogger.Printf("%+v", err)
		}
	}
}
