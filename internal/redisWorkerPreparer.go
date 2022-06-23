package internal

type RedisWorkerPreparer struct {
	worker *RedisWorker
}

func NewRedisWorkerPreparer(worker *RedisWorker) *RedisWorkerPreparer {
	return &RedisWorkerPreparer{
		worker: worker,
	}
}

func (p *RedisWorkerPreparer) RegisterErrorHandler(handler RedisErrorHandleProc) {
	p.worker.dispatcher.redisErrorHandler = handler
}

func (p *RedisWorkerPreparer) RegisterUnhandledMessageHandler(handler MessageHandler) {
	p.worker.dispatcher.unhandledMessageHandler = handler
}

func (p *RedisWorkerPreparer) RegisterStream(streamOffset StreamOffset) {
	p.worker.dispatcher.streams[streamOffset.Stream] = streamOffset
}

func (p *RedisWorkerPreparer) Router() Router {
	return p.worker.dispatcher.router
}
