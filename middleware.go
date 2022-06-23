package redis

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-redis/internal/middleware"
)

func UseErrorHandler(handler RedisErrorHandleProc) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.ErrorHandlerMiddleware{
		Handler: handler,
	}
}

func UseStreamGateway(streamGateway interface{}) host.Middleware {
	if streamGateway == nil {
		panic("argument 'topicGateway' cannot be nil")
	}

	return &middleware.StreamGatewayMiddleware{
		StreamGateway: streamGateway,
	}
}
