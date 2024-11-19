package internal

type ContextHelper struct{}

func (ContextHelper) ExtractReplyCode(ctx *Context) ReplyCode {
	if ctx != nil {
		reply, ok := ctx.values[__CONTEXT_REPLY_KEY]
		if ok {
			v, ok := reply.(ReplyCode)
			if ok {
				return v
			}
			return INVALID
		}
	}
	return UNSET
}

func (ContextHelper) InjectReplyCode(ctx *Context, reply ReplyCode) {
	ctx.values[__CONTEXT_REPLY_KEY] = reply
}

func (ContextHelper) InjectReplyCodeSafe(ctx *Context, reply ReplyCode) {
	_, ok := ctx.values[__CONTEXT_REPLY_KEY]
	if !ok {
		ctx.values[__CONTEXT_REPLY_KEY] = reply
	}
}
