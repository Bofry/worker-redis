package internal

import (
	"log"
	"os"
	"reflect"

	redis "github.com/Bofry/lib-redis-stream"
)

const (
	LOGGER_PREFIX string = "[worker-redis] "
)

var (
	RedisWorkerServiceInstance = new(RedisWorkerService)

	typeOfHost = reflect.TypeOf(RedisWorker{})

	logger *log.Logger = log.New(os.Stdout, LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	UniversalOptions = redis.UniversalOptions
	UniversalClient  = redis.UniversalClient
	XMessage         = redis.XMessage
	XStream          = redis.XStream

	StreamOffset   = redis.StreamOffset
	ConsumeContext = redis.ConsumeContext

	MessageHandler interface {
		ProcessMessage(ctx *ConsumeContext, stream string, message *XMessage)
	}
)

type (
	RedisErrorHandleProc = redis.RedisErrorHandleProc
	MessageHandleProc    = redis.MessageHandleProc
)
