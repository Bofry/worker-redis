package redis

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-redis/internal/middleware"
)

func UseErrorHandler(handler ErrorHandler) host.Middleware {
	if handler == nil {
		panic("argument 'handler' cannot be nil")
	}

	return &middleware.ErrorHandlerMiddleware{
		Handler: handler,
	}
}

func UseLogging(services ...LoggingService) host.Middleware {
	if len(services) == 0 {
		return &middleware.LoggingMiddleware{
			LoggingService: middleware.NoopLoggingServiceSingleton,
		}
	}

	return &middleware.LoggingMiddleware{
		LoggingService: middleware.NewCompositeLoggingService(services...),
	}
}

func UseMessageManager(messageManager interface{}) host.Middleware {
	if messageManager == nil {
		panic("argument 'messageManager' cannot be nil")
	}

	return &middleware.MessageManagerMiddleware{
		MessageManager: messageManager,
	}
}

func UseMessageObserverManager(messageObserverManager interface{}) host.Middleware {
	if messageObserverManager == nil {
		panic("argument 'messageObserverManager' cannot be nil")
	}

	return &middleware.MessageObserverManagerMiddleware{
		MessageObserverManager: messageObserverManager,
	}
}

func UseTracing(enabled bool) host.Middleware {
	return &middleware.TracingMiddleware{
		Enabled: enabled,
	}
}
