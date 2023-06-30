package internal

import (
	"context"
	"log"
)

type ContextBuilder struct {
	ctx *Context
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		ctx: new(Context),
	}
}

func (b *ContextBuilder) ConsumerGroup(v string) *ContextBuilder {
	b.ctx.ConsumerGroup = v
	return b
}

func (b *ContextBuilder) ConsumerName(v string) *ContextBuilder {
	b.ctx.ConsumerName = v
	return b
}

func (b *ContextBuilder) Context(v context.Context) *ContextBuilder {
	b.ctx.context = v
	return b
}

func (b *ContextBuilder) Logger(v *log.Logger) *ContextBuilder {
	b.ctx.logger = v
	return b
}

func (b *ContextBuilder) InvalidMessageHandler(v MessageHandler) *ContextBuilder {
	b.ctx.invalidMessageHandler = v
	return b
}

func (b *ContextBuilder) Build() *Context {
	return &Context{
		ConsumerGroup:         b.ctx.ConsumerGroup,
		ConsumerName:          b.ctx.ConsumerName,
		context:               b.ctx.context,
		logger:                b.ctx.logger,
		invalidMessageHandler: b.ctx.invalidMessageHandler,
	}
}
