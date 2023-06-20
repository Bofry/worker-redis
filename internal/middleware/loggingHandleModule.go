package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/worker-redis/internal"
)

var _ internal.MessageHandleModule = new(LoggingHandleModule)

type LoggingHandleModule struct {
	successor      internal.MessageHandleModule
	loggingService LoggingService
}

// CanSetSuccessor implements internal.MessageHandleModule.
func (*LoggingHandleModule) CanSetSuccessor() bool {
	return true
}

// SetSuccessor implements internal.MessageHandleModule.
func (m *LoggingHandleModule) SetSuccessor(successor internal.MessageHandleModule) {
	m.successor = successor
}

// ProcessMessage implements internal.MessageHandleModule.
func (m *LoggingHandleModule) ProcessMessage(ctx *internal.Context, message *redis.Message, state internal.ProcessingState, recover *internal.Recover) {
	if m.successor != nil {
		evidence := EventEvidence{
			traceID: state.Span.TraceID(),
			spanID:  state.Span.SpanID(),
			stream:  state.Stream,
		}

		eventLog := m.loggingService.CreateEventLog(evidence)
		eventLog.OnProcessMessage(message)

		recover.
			Defer(func(err interface{}) {
				if err != nil {
					defer func() {
						eventLog.OnError(message, err, debug.Stack())
						eventLog.Flush()
					}()

					// NOTE: we should not handle error here, due to the underlying RequestHandler
					// will handle it.
				} else {
					var (
						reply internal.ReplyCode = internal.GlobalContextHelper.ExtractReplyCode(ctx)
					)
					defer eventLog.Flush()

					eventLog.OnProcessMessageComplete(message, reply)
				}
			}).
			Do(func(f internal.Finalizer) {
				m.successor.ProcessMessage(ctx, message, state, recover)
			})
	}
}

// OnInitComplete implements internal.MessageHandleModule.
func (*LoggingHandleModule) OnInitComplete() {
	// ignored
}

// OnStart implements internal.MessageHandleModule.
func (*LoggingHandleModule) OnStart(ctx context.Context) error {
	// do nothing
	// return nil
	return fmt.Errorf("some errors occurred")
}

// OnStop implements internal.MessageHandleModule.
func (*LoggingHandleModule) OnStop(ctx context.Context) error {
	// do nothing
	return nil
}
