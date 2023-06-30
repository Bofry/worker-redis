package internal

import (
	"github.com/Bofry/trace"
)

var _ MessageHandler = new(StdInvalidMessageHandler)

type StdInvalidMessageHandler struct {
	invalidMessageHandler MessageHandler

	messageObserverService *MessageObserverService
}

// ProcessMessage implements MessageHandler.
func (h *StdInvalidMessageHandler) ProcessMessage(ctx *Context, message *Message) {
	GlobalContextHelper.InjectReplyCode(ctx, ABORT)

	if ctx.invalidMessageHandler != nil {
		var (
			sp     = trace.SpanFromContext(ctx)
			cloned = ctx.clone()
		)

		cloned.context = sp.Context()
		cloned.invalidMessageHandler = RestrictedForwardMessageHandler(RestrictedForwardMessage_Recursive)

		// prevent invalid message trigger observer methods
		h.messageObserverService.UnregisterAllMessageObservers(message)

		h.invalidMessageHandler.ProcessMessage(cloned, message)
	}
}
