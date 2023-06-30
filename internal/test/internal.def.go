package test

import (
	"log"
	"time"

	redis "github.com/Bofry/worker-redis"
)

var (
	defaultLogger *log.Logger = log.New(log.Writer(), "[worker-redis-test] ", log.LstdFlags|log.Lmsgprefix|log.LUTC)
)

type (
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

		// jaeger
		JaegerTraceUrl string `yaml:"jaegerTraceUrl"`
		JaegerQueryUrl string `yaml:"jaegerQueryUrl"`
	}
)

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
	h.AllowCreateGroup = true
}

func (h *Host) OnError(err error) (disposed bool) {
	return false
}
