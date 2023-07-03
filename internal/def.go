package internal

import (
	"context"
	"log"
	"os"
	"reflect"
	"sync/atomic"

	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

const (
	LOGGER_PREFIX string = "[worker-redis] "

	__CONTEXT_REPLY_KEY ctxReplyKeyType = 0
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
	typeOfMessageObserverAffair = reflect.TypeOf((*MessageObserverAffair)(nil)).Elem()
	defaultTracerProvider       = createNoopTracerProvider()
	defaultTextMapPropagator    = createNoopTextMapPropagator()

	globalTracerManager = defaultTracerManager()

	GlobalContextHelper             ContextHelper         = ContextHelper{}
	GlobalRestrictedMessageDelegate redis.MessageDelegate = RestrictedMessageDelegate(0)
	GlobalMessageHelper             MessageHelper         = MessageHelper{}

	RedisWorkerModuleInstance = RedisWorkerModule{}

	RedisWorkerLogger *log.Logger = log.New(os.Stdout, LOGGER_PREFIX, log.LstdFlags|log.Lmsgprefix)
)

type (
	ctxReplyKeyType int

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

func defaultTracerManager() *atomic.Value {
	v := &atomic.Value{}
	v.Store(tracerManagerHolder{
		v: NewTraceManager(),
	})
	return v
}

func GetTracerManager() *TracerManager {
	return globalTracerManager.Load().(tracerManagerHolder).v
}

func SetTracerManager(v *TracerManager) {
	current := GetTracerManager()
	if current != v {
		globalTracerManager.Store(tracerManagerHolder{
			v: v,
		})
	}
}
