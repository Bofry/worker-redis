package middleware

import (
	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/worker-redis/internal"
)

var _ EventLog = CompositeEventLog{}

type CompositeEventLog struct {
	EventLogs []EventLog
}

// Flush implements EventLog.
func (l CompositeEventLog) Flush() {
	for _, log := range l.EventLogs {
		log.Flush()
	}
}

// OnError implements EventLog.
func (l CompositeEventLog) OnError(message *redis.Message, err interface{}, stackTrace []byte) {
	for _, log := range l.EventLogs {
		log.OnError(message, err, stackTrace)
	}
}

// OnProcessMessage implements EventLog.
func (l CompositeEventLog) OnProcessMessage(message *redis.Message) {
	for _, log := range l.EventLogs {
		log.OnProcessMessage(message)
	}
}

// OnProcessMessageComplete implements EventLog.
func (l CompositeEventLog) OnProcessMessageComplete(message *redis.Message, reply internal.ReplyCode) {
	for _, log := range l.EventLogs {
		log.OnProcessMessageComplete(message, reply)
	}
}
