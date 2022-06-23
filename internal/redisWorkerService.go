package internal

import (
	"reflect"

	"github.com/Bofry/host"
)

var _ host.HostService = new(RedisWorkerService)

type RedisWorkerService struct{}

func (p *RedisWorkerService) Init(h host.Host, ctx *host.AppContext) {
	if v, ok := h.(*RedisWorker); ok {
		v.preInit()
	}
}

func (p *RedisWorkerService) InitComplete(h host.Host, ctx *host.AppContext) {
	if v, ok := h.(*RedisWorker); ok {
		v.init()
	}
}

func (p *RedisWorkerService) GetHostType() reflect.Type {
	return typeOfHost
}
