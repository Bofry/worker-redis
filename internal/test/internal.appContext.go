package standardtest

import (
	"fmt"
	"time"

	redis "github.com/Bofry/worker-redis"
)

type (
	App struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider
	}

	Host redis.Worker

	Config struct {
		// redis
		RedisAddresses           []string      `env:"*REDIS_SERVER"        yaml:"-"`
		RedisConsumerGroup       string        `env:"-"                    yaml:"RedisConsumerGroup"`
		RedisConsumerName        string        `env:"-"                    yaml:"RedisConsumerName"`
		RedisMaxInFlight         int64         `env:"-"                    yaml:"RedisMaxInFlight"`
		RedisMaxPollingTimeout   time.Duration `env:"-"                    yaml:"RedisMaxPollingTimeout"`
		RedisClaimMinIdleTime    time.Duration `env:"-"                    yaml:"RedisClaimMinIdleTime"`
		RedisIdlingTimeout       time.Duration `env:"-"                    yaml:"RedisIdlingTimeout"`
		RedisClaimSensitivity    int           `env:"-"                    yaml:"RedisClaimSensitivity"`
		RedisClaimOccurrenceRate int32         `env:"-"                    yaml:"RedisClaimOccurrenceRate"`
	}

	ServiceProvider struct {
		ResourceName string
	}

	StreamGateway struct {
		*GoTestStreamMessageHandler `stream:"gotestStream"`
		*UnhandledMessageHandler    `stream:"?"`
	}
)

func (provider *ServiceProvider) Init(conf *Config) {
	fmt.Println("ServiceProvider.Init()")
	provider.ResourceName = "demo resource"
}

func (h *Host) Init(conf *Config) {
	h.RedisOption = &redis.UniversalOptions{
		Addrs: conf.RedisAddresses,
	}
	h.ConsumerGroup = conf.RedisConsumerGroup
	h.ConsumerName = conf.RedisConsumerName
	h.MaxInFlight = conf.RedisMaxInFlight
	h.MaxPollingTimeout = conf.RedisMaxPollingTimeout
	h.ClaimMinIdleTime = conf.RedisClaimMinIdleTime
	h.IdlingTimeout = conf.RedisIdlingTimeout
	h.ClaimSensitivity = conf.RedisClaimSensitivity
	h.ClaimOccurrenceRate = conf.RedisClaimOccurrenceRate

}
