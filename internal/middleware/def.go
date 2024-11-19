package middleware

import (
	"log"
	"reflect"

	"github.com/Bofry/worker-redis/internal"
)

const (
	INVALID_MESSAGE_HANDLER_STREAM_SYMBOL = "?"

	TAG_STREAM                       = "stream"
	TAG_OFFSET                       = "offset"
	TAG_OPT_EXPAND_ENV               = "@ExpandEnv"
	TAG_OPT_MESSAGE_STATE_KEY_PREFIX = "@MessageStateKeyPrefix"
)

var (
	typeOfHost            = reflect.TypeOf(internal.RedisWorker{})
	typeOfMessageObserver = reflect.TypeOf((*internal.MessageObserver)(nil)).Elem()

	NoopLoggingServiceSingleton = NoopLoggingService{}
	NoopEventLogSingleton       = NoopEventLog(0)
)

type (
	ConfigureInvalidMessageHandleProc func(handler internal.MessageHandler)
	ConfigureStream                   func(stream internal.StreamOffset)

	LoggingService interface {
		CreateEventLog(ev EventEvidence) EventLog
		ConfigureLogger(l *log.Logger)
	}

	EventLog interface {
		OnError(message *internal.Message, err interface{}, stackTrace []byte)
		OnProcessMessage(message *internal.Message)
		OnProcessMessageComplete(message *internal.Message, reply internal.ReplyCode)
		Flush()
	}
)
