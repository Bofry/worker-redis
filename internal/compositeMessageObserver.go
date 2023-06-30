package internal

var _ MessageObserver = CompositeMessageObserver(nil)

type CompositeMessageObserver []MessageObserver

// OnAck implements MessageObserver.
func (o CompositeMessageObserver) OnAck(ctx *Context, message *Message) {
	cloned := message.Clone()
	cloned.Delegate = GlobalRestrictedMessageDelegate

	for _, handler := range o {
		handler.OnAck(ctx, cloned)
	}
}

// OnDel implements MessageObserver.
func (o CompositeMessageObserver) OnDel(ctx *Context, message *Message) {
	cloned := message.Clone()
	cloned.Delegate = GlobalRestrictedMessageDelegate

	for _, handler := range o {
		handler.OnDel(ctx, cloned)
	}
}
