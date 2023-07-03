package internal

import (
	"context"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Bofry/host"
	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

var _ host.Host = new(RedisWorker)

type RedisWorker struct {
	ConsumerGroup       string
	ConsumerName        string
	RedisOption         *redis.UniversalOptions
	MaxInFlight         int64
	MaxPollingTimeout   time.Duration
	ClaimMinIdleTime    time.Duration
	IdlingTimeout       time.Duration // 若沒有任何訊息時等待多久
	ClaimSensitivity    int           // Read 時取得的訊息數小於等於 n 的話, 執行 Claim
	ClaimOccurrenceRate int32         // Read 每執行 n 次後 執行 Claim 1 次
	AllowCreateGroup    bool          // 自動註冊 consumer group

	consumer *redis.Consumer

	logger *log.Logger

	messageDispatcher *MessageDispatcher
	messageManager    interface{}

	messageHandleService   *MessageHandleService
	messageTracerService   *MessageTracerService
	messageObserverService *MessageObserverService

	tracerManager *TracerManager

	onErrorEventHandler host.HostOnErrorEventHandler

	wg          sync.WaitGroup
	mutex       sync.Mutex
	initialized bool
	running     bool
	disposed    bool
}

func (w *RedisWorker) Start(ctx context.Context) {
	if w.disposed {
		RedisWorkerLogger.Panic("the Worker has been disposed")
	}
	if !w.initialized {
		RedisWorkerLogger.Panic("the Worker havn't be initialized yet")
	}
	if w.running {
		return
	}

	var err error
	w.mutex.Lock()
	defer func() {
		if err != nil {
			w.running = false
			w.disposed = true
		}
		w.mutex.Unlock()
	}()

	w.running = true
	w.messageDispatcher.start(ctx)

	var (
		streams       = w.messageDispatcher.Streams()
		streamOffsets = w.messageDispatcher.StreamOffsets()
	)

	RedisWorkerLogger.Printf("name [%s] group [%s] listening DB [%d] streams [%s] on address %s\n",
		w.ConsumerName,
		w.ConsumerGroup,
		w.RedisOption.DB,
		strings.Join(streams, ","),
		w.RedisOption.Addrs)

	if len(streamOffsets) > 0 {
		c := w.consumer
		err = w.registerGroup(streamOffsets)
		if err != nil {
			RedisWorkerLogger.Panic(err)
		}
		err = w.messageDispatcher.subscribe(c)
		if err != nil {
			RedisWorkerLogger.Panic(err)
		}
	}
}

func (w *RedisWorker) Stop(ctx context.Context) error {
	RedisWorkerLogger.Printf("%% Stopping\n")

	w.mutex.Lock()
	defer func() {
		w.running = false
		w.disposed = true
		w.mutex.Unlock()

		w.messageDispatcher.stop(ctx)

		RedisWorkerLogger.Printf("%% Stopped\n")
	}()

	w.consumer.Close()
	w.wg.Wait()
	return nil
}

func (w *RedisWorker) Logger() *log.Logger {
	return w.logger
}

func (w *RedisWorker) alloc() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.tracerManager = NewTraceManager()
	w.messageHandleService = NewMessageHandleService()
	w.messageTracerService = &MessageTracerService{
		TracerManager: w.tracerManager,
	}
	w.messageObserverService = &MessageObserverService{
		MessageObservers: make(map[reflect.Type]MessageObserver),
	}

	w.messageDispatcher = &MessageDispatcher{
		MessageHandleService:   w.messageHandleService,
		MessageTracerService:   w.messageTracerService,
		MessageObserverService: w.messageObserverService,
		Router:                 make(Router),
		StreamSet:              make(map[string]StreamOffset),
		OnHostErrorProc:        w.onHostError,
	}

	// register TracerManager
	SetTracerManager(w.tracerManager)
}

func (w *RedisWorker) init() {
	if w.initialized {
		return
	}

	w.mutex.Lock()
	defer func() {
		w.initialized = true
		w.mutex.Unlock()
	}()

	w.messageTracerService.init(w.messageManager)
	w.messageObserverService.init(w.messageManager)
	w.messageDispatcher.init()
	w.configConsumer()
}

func (w *RedisWorker) registerGroup(offsets []StreamOffset) error {
	if !w.AllowCreateGroup {
		return nil
	}

	var (
		admin *redis.AdminClient
		err   error
	)
	w.wg.Add(1)
	defer w.wg.Done()

	admin, err = redis.NewAdminClient(w.RedisOption)
	if err != nil {
		return err
	}
	defer admin.Close()

	// XGROUP CREATE AND MKSTREAM
	for _, readOffset := range offsets {
		if len(readOffset.Offset) == 0 {
			readOffset.Offset = redis.StreamLastDeliveredID
		}

		_, err := admin.CreateConsumerGroupAndStream(readOffset.Stream, w.ConsumerGroup, readOffset.Offset)
		if err != nil {
			if !(isRedisBusyGroupError(err)) {
				return err
			}
		}
	}
	return nil
}

func (w *RedisWorker) configConsumer() {
	instance := &redis.Consumer{
		Group:               w.ConsumerGroup,
		Name:                w.ConsumerName,
		RedisOption:         w.RedisOption,
		MaxInFlight:         w.MaxInFlight,
		MaxPollingTimeout:   w.MaxPollingTimeout,
		ClaimMinIdleTime:    w.ClaimMinIdleTime,
		IdlingTimeout:       w.IdlingTimeout,
		ClaimSensitivity:    w.ClaimSensitivity,
		ClaimOccurrenceRate: w.ClaimOccurrenceRate,
		MessageHandler:      w.receiveMessage,
		RedisErrorHandler:   w.onHostError,
		Logger:              w.logger,
	}

	w.consumer = instance
}

func (w *RedisWorker) receiveMessage(message *Message) {
	ctx := &Context{
		ConsumerGroup:         w.ConsumerGroup,
		ConsumerName:          w.ConsumerName,
		logger:                w.logger,
		invalidMessageHandler: nil, // be determined by MessageDispatcher
	}

	// configure nsq.MessageDelegate
	delegate := NewContextMessageDelegate(ctx)
	delegate.configure(message)

	w.messageDispatcher.ProcessMessage(ctx, message)
}

func (w *RedisWorker) onHostError(err error) (disposed bool) {
	if w.onErrorEventHandler != nil {
		return w.onErrorEventHandler.OnError(err)
	}
	return false
}

func (w *RedisWorker) setTextMapPropagator(propagator propagation.TextMapPropagator) {
	w.messageTracerService.textMapPropagator = propagator
}

func (w *RedisWorker) setTracerProvider(provider *trace.SeverityTracerProvider) {
	w.messageTracerService.tracerProvider = provider
}

func (w *RedisWorker) setLogger(l *log.Logger) {
	w.logger = l
}
