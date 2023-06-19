package redis

import (
	"github.com/Bofry/host"
	"github.com/Bofry/worker-redis/internal"
)

func Startup(app interface{}) *host.Starter {
	var (
		starter = host.Startup(app)
	)

	host.RegisterHostModule(starter, internal.RedisWorkerModuleInstance)

	return starter
}
