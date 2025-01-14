package internal

var _ MessageHandler = new(MessageHandleProc)

type MessageHandleProc func(ctx *Context, message *Message)

// ProcessMessage implements MessageHandler.
func (proc MessageHandleProc) ProcessMessage(ctx *Context, message *Message) {
	proc(ctx, message)
}

// //////////////////////////////////////////
var (
	_ MessageHandleProc = StopRecursiveForwardMessageHandler
)

func StopRecursiveForwardMessageHandler(ctx *Context, msg *Message) {
	ctx.logger.Fatal("invalid forward; it might be recursive forward message to InvalidMessageHandler")
}
