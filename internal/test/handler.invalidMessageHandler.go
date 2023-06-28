package test

import (
	"fmt"

	redis "github.com/Bofry/worker-redis"
)

type InvalidMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *InvalidMessageHandler) Init() {
	fmt.Println("InvalidMessageHandler.Init()")
}

func (h *InvalidMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	ctx.Logger().Printf("Invalid Message on %s: %v\n", message.Stream, message.XMessage)
}
