package internal

import (
	"reflect"

	redis "github.com/Bofry/lib-redis-stream"
)

type RedisWorkerRegistrar struct {
	worker *RedisWorker
}

func NewRedisWorkerRegistrar(worker *RedisWorker) *RedisWorkerRegistrar {
	return &RedisWorkerRegistrar{
		worker: worker,
	}
}

func (r *RedisWorkerRegistrar) RegisterMessageHandleModule(module MessageHandleModule) {
	r.worker.messageHandleService.Register(module)
}

func (r *RedisWorkerRegistrar) EnableTracer(enabled bool) {
	r.worker.messageTracerService.Enabled = enabled
}

func (r *RedisWorkerRegistrar) SetErrorHandler(handler ErrorHandler) {
	r.worker.messageDispatcher.ErrorHandler = handler
}

func (r *RedisWorkerRegistrar) SetInvalidMessageHandler(handler MessageHandler) {
	r.worker.messageDispatcher.InvalidMessageHandler = handler
}

func (r *RedisWorkerRegistrar) SetMessageManager(messageManager interface{}) {
	r.worker.messageManager = messageManager
}

func (r *RedisWorkerRegistrar) RegisterStream(streamOffset redis.StreamOffset) {
	r.worker.messageDispatcher.StreamSet[streamOffset.Stream] = streamOffset
}

func (r *RedisWorkerRegistrar) AddRouter(stream string, handler MessageHandler, handlerComponentID string) {
	r.worker.messageDispatcher.Router.Add(stream, handler, handlerComponentID)
}

func (r *RedisWorkerRegistrar) RegisterMessageObserver(v MessageObserver) {
	t := reflect.TypeOf(v)
	r.worker.messageObserverService.MessageObservers[t] = v
}
