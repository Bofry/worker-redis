package test

import (
	"fmt"

	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
)

var (
	_ redis.MessageHandler = new(InvalidMessageHandler)
)

type InvalidMessageHandler struct {
	ServiceProvider *ServiceProvider

	counter *TestStreamMessageCounter
}

func (h *InvalidMessageHandler) Init() {
	fmt.Println("InvalidMessageHandler.Init()")

	h.counter = h.ServiceProvider.TestStreamMessageCounter
}

// ProcessMessage implements internal.MessageHandler.
func (h *InvalidMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	sp := trace.SpanFromContext(ctx)

	sp.Info("InvalidMessage %+v", string(message.ID))

	ctx.Logger().Printf("Invalid Message on %s: %v\n", message.Stream, message.XMessage)

	h.counter.IncreaseInvalidMessageCount(sp.Context())
}
