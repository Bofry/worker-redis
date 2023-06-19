package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/worker-redis/internal"
)

var _ host.Middleware = new(LoggingMiddleware)

type LoggingMiddleware struct {
	LoggingService LoggingService
}

// Init implements internal.Middleware
func (m *LoggingMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asRedisWorker(app.Host())
		registrar = NewRedisWorkerRegistrar(worker)
	)

	m.LoggingService.ConfigureLogger(worker.Logger())

	loggingHandleModule := &LoggingHandleModule{
		loggingService: m.LoggingService,
	}
	registrar.RegisterMessageHandleModule(loggingHandleModule)
}
