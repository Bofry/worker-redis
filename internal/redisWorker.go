package internal

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Bofry/host"
	redis "github.com/Bofry/lib-redis-stream"
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

	consumer *redis.Consumer

	dispatcher *RedisMessageDispatcher

	wg          sync.WaitGroup
	mutex       sync.Mutex
	initialized bool
	running     bool
	disposed    bool
}

func (w *RedisWorker) Start(ctx context.Context) {
	if w.disposed {
		logger.Panic("the Worker has been disposed")
	}
	if !w.initialized {
		logger.Panic("the Worker havn't be initialized yet")
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

	var (
		streams       = w.dispatcher.Streams()
		streamOffsets = w.dispatcher.StreamOffsets()
	)

	logger.Printf("name [%s] group [%s] listening DB [%d] streams [%s] on address %s\n",
		w.ConsumerName,
		w.ConsumerGroup,
		w.RedisOption.DB,
		strings.Join(streams, ","),
		w.RedisOption.Addrs)

	if len(streamOffsets) > 0 {
		c := w.consumer
		err = c.Subscribe(streamOffsets...)
		if err != nil {
			logger.Panic(err)
		}
	}
}

func (w *RedisWorker) Stop(ctx context.Context) error {
	logger.Printf("%% Stopping\n")
	defer func() {
		logger.Printf("%% Stopped\n")
	}()

	w.consumer.Close()
	return nil
}

func (w *RedisWorker) preInit() {
	w.dispatcher = NewRedisMessageDispatcher()
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

	w.configConsumer()
}

func (w *RedisWorker) configConsumer() {
	instance := &redis.Consumer{
		Group:                   w.ConsumerGroup,
		Name:                    w.ConsumerName,
		RedisOption:             w.RedisOption,
		MaxInFlight:             w.MaxInFlight,
		MaxPollingTimeout:       w.MaxPollingTimeout,
		ClaimMinIdleTime:        w.ClaimMinIdleTime,
		IdlingTimeout:           w.IdlingTimeout,
		ClaimSensitivity:        w.ClaimSensitivity,
		ClaimOccurrenceRate:     w.ClaimOccurrenceRate,
		ErrorHandler:            w.dispatcher.ProcessRedisError,
		MessageHandler:          w.dispatcher.ProcessMessage,
		UnhandledMessageHandler: w.dispatcher.ProcessUnhandledMessage,
	}

	w.consumer = instance
}
