package standardtest

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

func (h *UnhandledMessageHandler) ProcessMessage(ctx *redis.WorkerContext, stream string, message *redis.XMessage) {
	log.Printf("Unhandled Message on %s: %v\n", stream, message)
}
