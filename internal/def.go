package internal

import (
	"context"
	"log"
	"os"
	"reflect"

	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

const (
	LOGGER_PREFIX string = "[worker-redis] "

	__CONTEXT_REPLY_KEY ctxReplyKeyType = 0

	__UNDEFINED_TRACER_NAME     = "undefined"
	__INVALID_MESSAGE_SPAN_NAME = "InvalidMessage"

	__INVALID_MESSAGE_HANDLER_NAME = "InvalidMessageHandler"
)

const (
	UNSET ReplyCode = iota
	PASS
	FAIL
	ABORT

	__reply_code_minimum__ = UNSET
	__reply_code_maximum__ = ABORT

	INVALID ReplyCode = -1

	__reply_code_invalid_text__ = "invalid"
)

const (
	RestrictedForwardMessage_InvalidOperation int = 0
	RestrictedForwardMessage_Recursive        int = 1
)

var (
	typeOfHost                  = reflect.TypeOf(RedisWorker{})
	typeOfMessageHandler        = reflect.TypeOf((*MessageHandler)(nil)).Elem()
	typeOfMessageObserverAffair = reflect.TypeOf((*MessageObserverAffair)(nil)).Elem()
	defaultTracerProvider       = createNoopTracerProvider()
	defaultTextMapPropagator    = createNoopTextMapPropagator()
	defaultMessageDelegate      = NoopMessageDelegate(0)

	GlobalTracerManager             *TracerManager        // be register from NsqWorker
	GlobalContextHelper             ContextHelper         = ContextHelper{}
	GlobalRestrictedMessageDelegate redis.MessageDelegate = RestrictedMessageDelegate(0)
	GlobalNoopMessageDelegate       redis.MessageDelegate = NoopMessageDelegate(0)
	GlobalMessageDelegateHelper     MessageDelegateHelper = MessageDelegateHelper{}

	RedisWorkerModuleInstance = RedisWorkerModule{}

	RedisWorkerLogger *log.Logger = log.New(os.Stdout, LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	ctxReplyKeyType int

	StatusCode = ReplyCode

	tracerManagerHolder struct {
		v *TracerManager
	}

	UniversalOptions = redis.UniversalOptions
	UniversalClient  = redis.UniversalClient
	XMessage         = redis.XMessage
	XStream          = redis.XStream
	Message          = redis.Message
	StreamOffset     = redis.StreamOffset
	StreamOffsetInfo = redis.StreamOffsetInfo

	MessageHandleModule interface {
		CanSetSuccessor() bool
		SetSuccessor(successor MessageHandleModule)
		ProcessMessage(ctx *Context, message *Message, state ProcessingState, recover *Recover)
		OnInitComplete()
		OnStart(ctx context.Context) error
		OnStop(ctx context.Context) error
	}

	MessageObserver interface {
		OnAck(ctx *Context, message *Message)
		OnDel(ctx *Context, message *Message)
		Type() reflect.Type
	}

	MessageObserverAffair interface {
		MessageObserverTypes() []reflect.Type
	}

	MessageHandler interface {
		ProcessMessage(ctx *Context, message *Message)
	}

	ErrorHandler func(ctx *Context, message *Message, err interface{})

	OnHostErrorHandler func(err error) (disposed bool)
)

func createNoopTracerProvider() *trace.SeverityTracerProvider {
	tp, err := trace.NoopProvider()
	if err != nil {
		RedisWorkerLogger.Fatalf("cannot create NoopProvider: %v", err)
	}
	return tp
}

func createNoopTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator()
}
