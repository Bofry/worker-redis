package test

import (
	"fmt"

	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
)

type InvalidMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *InvalidMessageHandler) Init() {
	fmt.Println("InvalidMessageHandler.Init()")
}

func (h *InvalidMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	sp := trace.SpanFromContext(ctx)

	sp.Info("InvalidMessage %+v", string(message.ID))

	ctx.Logger().Printf("Invalid Message on %s: %v\n", message.Stream, message.XMessage)
}
