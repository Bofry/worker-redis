package internal

type RedisMessageDispatcher struct {
	redisErrorHandler       RedisErrorHandleProc
	messageHandler          MessageHandleProc
	unhandledMessageHandler MessageHandler
	router                  Router
	streams                 map[string]StreamOffset
}

func NewRedisMessageDispatcher() *RedisMessageDispatcher {
	return &RedisMessageDispatcher{
		router:  make(Router),
		streams: make(map[string]StreamOffset),
	}
}

func (d *RedisMessageDispatcher) StreamOffsets() []StreamOffset {
	var (
		streams = d.streams
	)

	offsets := make([]StreamOffset, 0, len(streams))
	for _, v := range streams {
		offsets = append(offsets, v)
	}
	return offsets
}

func (d *RedisMessageDispatcher) Streams() []string {
	var (
		router = d.router
	)

	if router != nil {
		keys := make([]string, 0, len(router))
		for k := range router {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *RedisMessageDispatcher) ProcessMessage(ctx *ConsumeContext, stream string, message *XMessage) {

	// TODO: handle error
	// defer func() {
	// 	err := recover()
	// 	if err != nil {
	// 		d.ProcessRedisError(err)
	// 	}
	// }()

	handler := d.router.Get(stream)
	if handler != nil {
		handler.ProcessMessage(ctx, stream, message)
	} else {
		ctx.ForwardUnhandledMessage(stream, message)
	}
}

func (d *RedisMessageDispatcher) ProcessUnhandledMessage(ctx *ConsumeContext, stream string, message *XMessage) {
	if d.unhandledMessageHandler != nil {
		d.unhandledMessageHandler.ProcessMessage(ctx, stream, message)
	}
}

func (d *RedisMessageDispatcher) ProcessRedisError(err error) (disposed bool) {
	if d.redisErrorHandler != nil {
		return d.redisErrorHandler(err)
	}
	return false
}
