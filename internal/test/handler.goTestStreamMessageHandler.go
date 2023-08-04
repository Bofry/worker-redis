package test

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
	"github.com/Bofry/worker-redis/tracing"
)

var (
	_ redis.MessageHandler        = new(GoTestStreamMessageHandler)
	_ redis.MessageObserverAffair = new(GoTestStreamMessageHandler)
)

type GoTestStreamMessageHandler struct {
	ServiceProvider *ServiceProvider

	counter *GoTestStreamMessageCounter
}

func (h *GoTestStreamMessageHandler) Init() {
	fmt.Println("GoTestStreamMessageHandler.Init()")

	h.counter = new(GoTestStreamMessageCounter)
}

// ProcessMessage implements internal.MessageHandler.
func (h *GoTestStreamMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	ctx.Logger().Printf("Message on %s: %v\n", message.Stream, message.XMessage)

	sp := trace.SpanFromContext(ctx)
	sp.Argv(fmt.Sprintf("%+v", message.Values))

	if message.Stream == "gotestStream2" {
		h.doSomething(sp.Context())
		ctx.InvalidMessage(message)
		return
	}
	h.counter.increase(sp.Context())
	message.Ack()
}

// MessageObserverTypes implements internal.MessageObserverAffair.
func (*GoTestStreamMessageHandler) MessageObserverTypes() []reflect.Type {
	return []reflect.Type{
		MessageObserverManager.GoTestStreamMessageObserver.Type(),
	}
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
