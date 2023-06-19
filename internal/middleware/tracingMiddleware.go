package middleware

import (
	"github.com/Bofry/host"
	. "github.com/Bofry/worker-redis/internal"
)

var _ host.Middleware = new(TracingMiddleware)

type TracingMiddleware struct {
	Enabled bool
}

// Init implements internal.Middleware.
func (m *TracingMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asRedisWorker(app.Host())
		registrar = NewRedisWorkerRegistrar(worker)
	)

	registrar.EnableTracer(m.Enabled)
}
