package internal

import "github.com/Bofry/trace"

var _ MessageObserver = CompositeMessageObserver(nil)

type CompositeMessageObserver []MessageObserver

// OnAck implements MessageObserver.
func (o CompositeMessageObserver) OnAck(ctx *Context, message *Message) {
	clonedMessage := message.Clone()
	clonedMessage.Delegate = GlobalRestrictedMessageDelegate

	var (
		sp        = trace.SpanFromContext(ctx)
		clonedCtx = ctx.clone()
	)
	clonedCtx.context = sp.Context()
	clonedCtx.invalidMessageHandler = RestrictedForwardMessageHandler(RestrictedForwardMessage_InvalidOperation)

	for _, handler := range o {
		handler.OnAck(clonedCtx, clonedMessage)
	}
}

// OnDel implements MessageObserver.
func (o CompositeMessageObserver) OnDel(ctx *Context, message *Message) {
	clonedMessage := message.Clone()
	clonedMessage.Delegate = GlobalRestrictedMessageDelegate

	var (
		sp        = trace.SpanFromContext(ctx)
		clonedCtx = ctx.clone()
	)
	clonedCtx.context = sp.Context()
	clonedCtx.invalidMessageHandler = RestrictedForwardMessageHandler(RestrictedForwardMessage_InvalidOperation)

	for _, handler := range o {
		handler.OnDel(clonedCtx, clonedMessage)
	}
}
