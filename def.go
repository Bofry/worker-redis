package redis

import (
	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/worker-redis/internal"
	"github.com/Bofry/worker-redis/internal/middleware"
)

const (
	StreamAsteriskID           = redis.StreamAsteriskID
	StreamLastDeliveredID      = redis.StreamLastDeliveredID
	StreamZeroID               = redis.StreamZeroID
	StreamZeroOffset           = redis.StreamZeroOffset
	StreamNeverDeliveredOffset = redis.StreamNeverDeliveredOffset
)

type (
	ProducerConfig   = redis.ProducerConfig
	Producer         = redis.Producer
	UniversalOptions = redis.UniversalOptions
	Message          = redis.Message
	MessageContent   = redis.MessageContent

	AdminClient     = redis.AdminClient
	Forwarder       = redis.Forwarder
	ForwarderRunner = redis.ForwarderRunner

	EventEvidence  = middleware.EventEvidence
	LoggingService = middleware.LoggingService
	EventLog       = middleware.EventLog

	MessageObserver       = internal.MessageObserver
	MessageObserverAffair = internal.MessageObserverAffair

	MessageHandler = internal.MessageHandler
	Worker         = internal.RedisWorker
	Context        = internal.Context
	ReplyCode      = internal.ReplyCode

	ErrorHandler = internal.ErrorHandler
)

func NewAdminClient(opt *UniversalOptions) (*AdminClient, error) {
	return redis.NewAdminClient(opt)
}

func NewForwarder(config *ProducerConfig) (*Forwarder, error) {
	return redis.NewForwarder(config)
}
