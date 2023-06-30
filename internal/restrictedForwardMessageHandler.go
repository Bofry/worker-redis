package internal

var (
	_RestrictedForwardMessageErrorMap = map[int]string{
		RestrictedForwardMessage_InvalidOperation: "restrict forward error; invalid operstion",
		RestrictedForwardMessage_Recursive:        "restrict forward error; recursive forward message to handler",
	}
	_DefaultRestrictedForwardMessageError = _RestrictedForwardMessageErrorMap[RestrictedForwardMessage_InvalidOperation]
)

var _ MessageHandler = RestrictedForwardMessageHandler(0)

type RestrictedForwardMessageHandler int

// ProcessMessage implements MessageHandler.
func (h RestrictedForwardMessageHandler) ProcessMessage(ctx *Context, message *Message) {
	var code = int(h)
	if msg, ok := _RestrictedForwardMessageErrorMap[code]; ok {
		panic(RestrictedOperationError(msg))
	}
	panic(RestrictedOperationError(_DefaultRestrictedForwardMessageError))
}
