package redis

import (
	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/worker-redis/internal"
)

const (
	StreamAsteriskID           = redis.StreamAsteriskID
	StreamLastDeliveredID      = redis.StreamLastDeliveredID
	StreamZeroID               = redis.StreamZeroID
	StreamZeroOffset           = redis.StreamZeroOffset
	StreamNeverDeliveredOffset = redis.StreamNeverDeliveredOffset
)

type (
	UniversalOptions = redis.UniversalOptions
	XMessage         = redis.XMessage
	XStream          = redis.XStream

	AdminClient     = redis.AdminClient
	Forwarder       = redis.Forwarder
	ForwarderRunner = redis.ForwarderRunner
	WorkerContext   = redis.ConsumeContext

	MessageHandler = internal.MessageHandler
	Worker         = internal.RedisWorker
)

type (
	RedisErrorHandleProc = internal.RedisErrorHandleProc
)
