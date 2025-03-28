package test

import (
	"context"
	"fmt"

	redislib "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
	"github.com/Bofry/worker-redis/tracing"
)

var (
	_ redis.MessageHandler        = new(GoTestStreamFilterMessageHandler)
	_ redis.MessageFilterAffinity = new(GoTestStreamFilterMessageHandler)
)

type GoTestStreamFilterMessageHandler struct {
	ServiceProvider *ServiceProvider

	counter *TestStreamMessageCounter
}

// Filter implements internal.MessageFilterAffinity.
func (h *GoTestStreamFilterMessageHandler) Filter(message *redis.Message) bool {
	fmt.Printf("stream: %s, values: %+v\n", message.Stream, message.Values)
	if len(message.Values) > 0 {
		if message.Values["name"] == "luffy" {
			return false
		}
	}
	return true
}

func (h *GoTestStreamFilterMessageHandler) Init() {
	fmt.Println("GoTestStreamFilterMessageHandler.Init()")

	h.counter = h.ServiceProvider.TestStreamMessageCounter
}

// ProcessMessage implements internal.MessageHandler.
func (h *GoTestStreamFilterMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
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

func (h *GoTestStreamFilterMessageHandler) doSomething(ctx context.Context) {
	tr := tracing.GetTracer(h)
	sp := tr.Start(ctx, "doSomething()")
	defer sp.End()
}
