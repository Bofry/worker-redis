package test

import (
	"log"

	redis "github.com/Bofry/worker-redis"
)

type GoTestStreamMessageHandler struct {
	ServiceProvider *ServiceProvider
}

func (h *GoTestStreamMessageHandler) Init() {
	log.Printf("GoTestStreamMessageHandler.Init()")
}

func (h *GoTestStreamMessageHandler) ProcessMessage(ctx *redis.WorkerContext, stream string, message *redis.XMessage) {
	log.Printf("Message on %s: %v\n", stream, message)
	ctx.Ack(stream, message.ID)
	ctx.ForwardUnhandledMessage(stream, message)
}
