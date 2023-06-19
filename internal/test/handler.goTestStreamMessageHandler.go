package test

import (
	"context"
	"fmt"
	"log"

	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
	"github.com/Bofry/worker-redis/tracing"
)

type GoTestStreamMessageHandler struct {
	ServiceProvider *ServiceProvider

	counter *GoTestStreamMessageCounter
}

func (h *GoTestStreamMessageHandler) Init() {
	log.Printf("GoTestStreamMessageHandler.Init()")

	h.counter = new(GoTestStreamMessageCounter)
}

func (h *GoTestStreamMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	log.Printf("Message on %s: %v\n", message.Stream, message.XMessage)

	sp := trace.SpanFromContext(ctx)
	sp.Argv(fmt.Sprintf("%+v", message.Values))

	if message.Stream == "gotestStream2" {
		h.doSomething(sp.Context())
		ctx.ThrowInvalidMessageError(message)
		return
	}
	h.counter.increase(sp.Context())
	message.Ack()
}

func (h *GoTestStreamMessageHandler) doSomething(ctx context.Context) {
	tr := tracing.GetTracer(h)
	sp := tr.Start(ctx, "doSomething()")
	defer sp.End()
}

type GoTestStreamMessageCounter struct {
	count int
}

func (c *GoTestStreamMessageCounter) increase(ctx context.Context) int {
	tr := tracing.GetTracer(c)
	sp := tr.Start(ctx, "increase()")
	defer sp.End()

	c.count++
	return c.count
}
