package internal

import (
	"io"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/host/helper"
)

var _ host.HostModule = RedisWorkerModule{}

type RedisWorkerModule struct{}

// ConfigureLogger implements host.HostModule
func (RedisWorkerModule) ConfigureLogger(logflags int, w io.Writer) {
	RedisWorkerLogger.SetFlags(logflags)
	RedisWorkerLogger.SetOutput(w)
}

func (RedisWorkerModule) Init(h host.Host, app *host.AppModule) {
	if v, ok := h.(*RedisWorker); ok {
		v.alloc()
		v.setTracerProvider(app.TracerProvider())
		v.setTextMapPropagator(app.TextMapPropagator())
		v.setLogger(app.Logger())

		{
			host := helper.HostHelper(app)
			v.onErrorEventHandler = host.OnErrorEventHandler()
		}
	}
}

func (RedisWorkerModule) InitComplete(h host.Host, ctx *host.AppModule) {
	if v, ok := h.(*RedisWorker); ok {
		v.init()
	}
}

func (RedisWorkerModule) DescribeHostType() reflect.Type {
	return typeOfHost
}
