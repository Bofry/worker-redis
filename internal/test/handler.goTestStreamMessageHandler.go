package test

import (
	"context"
	"fmt"
	"reflect"

	redislib "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
	"github.com/Bofry/worker-redis/tracing"
)

var (
	_ redis.MessageHandler          = new(GoTestStreamMessageHandler)
	_ redis.MessageObserverAffinity = new(GoTestStreamMessageHandler)
)

type GoTestStreamMessageHandler struct {
	ServiceProvider *ServiceProvider

	counter *TestStreamMessageCounter
}

func (h *GoTestStreamMessageHandler) Init() {
	fmt.Println("GoTestStreamMessageHandler.Init()")

	h.counter = h.ServiceProvider.TestStreamMessageCounter
}

// ProcessMessage implements internal.MessageHandler.
func (h *GoTestStreamMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	ctx.Logger().Printf("Message on %s: %v\n", message.Stream, message.XMessage)

	sp := trace.SpanFromContext(ctx)
	sp.Argv(fmt.Sprintf("%+v", message.Values))

	h.counter.IncreaseMessageCount(sp.Context())

	if message.Stream == "gotestStream2" {
		h.doSomething(sp.Context())
		ctx.InvalidMessage(message)
		return
	}
	if message.Stream == "gotestStream3" {
		var state = make(map[string]interface{})
		message.Content().State.Visit(func(name string, value interface{}) {
			state[name] = value
		})
		fmt.Printf("%+v\n", state)
	}
	if message.Stream == "gotestStream4" {
		var state = make(map[string]interface{})
		message.Content(
			redislib.WithMessageStateKeyPrefix("mystate:"),
		).State.Visit(func(name string, value interface{}) {
			state[name] = value
		})
		fmt.Printf("%+v\n", state)
	}
	if message.Stream == "gotestStream5" {
		panic("something occurred")
	}

	h.counter.IncreaseSuccessMessageCount(sp.Context())
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
