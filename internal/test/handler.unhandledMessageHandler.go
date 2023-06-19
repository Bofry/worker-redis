package test

import (
	"log"

	redis "github.com/Bofry/worker-redis"
)

type UnhandledMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *UnhandledMessageHandler) Init() {
	log.Printf("UnhandledMessageHandler.Init()")
}

func (h *UnhandledMessageHandler) ProcessMessage(ctx *redis.Context, message *redis.Message) {
	log.Printf("Unhandled Message on %s: %v\n", message.Stream, message.XMessage)
}
