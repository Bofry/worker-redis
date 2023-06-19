package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/worker-redis/internal"
)

var _ host.Middleware = new(ErrorHandlerMiddleware)

type ErrorHandlerMiddleware struct {
	Handler ErrorHandler
}

func (m *ErrorHandlerMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asRedisWorker(app.Host())
		registrar = NewRedisWorkerRegistrar(worker)
	)

	registrar.SetErrorHandler(m.Handler)
}
