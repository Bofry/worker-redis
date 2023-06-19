package test

import (
	"log"

	redis "github.com/Bofry/worker-redis"
)

var _ redis.LoggingService = new(MultiLoggerService)

type MultiLoggerService struct {
	LoggingServices []redis.LoggingService
}

func (s *MultiLoggerService) CreateEventLog(ev redis.EventEvidence) redis.EventLog {
	var eventlogs []redis.EventLog
	for _, svc := range s.LoggingServices {
		eventlogs = append(eventlogs, svc.CreateEventLog(ev))
	}

	return MultiEventLog{
		EventLogs: eventlogs,
	}
}

func (s *MultiLoggerService) ConfigureLogger(l *log.Logger) {
	for _, svc := range s.LoggingServices {
		svc.ConfigureLogger(l)
	}
}

var _ redis.EventLog = MultiEventLog{}

type MultiEventLog struct {
	EventLogs []redis.EventLog
}

// Flush implements middleware.EventLog.
func (l MultiEventLog) Flush() {
	for _, log := range l.EventLogs {
		log.Flush()
	}
}

// OnError implements middleware.EventLog.
func (l MultiEventLog) OnError(message *redis.Message, err interface{}, stackTrace []byte) {
	for _, log := range l.EventLogs {
		log.OnError(message, err, stackTrace)
	}
}

// OnProcessMessageComplete implements middleware.EventLog.
func (l MultiEventLog) OnProcessMessageComplete(message *redis.Message, reply redis.ReplyCode) {
	for _, log := range l.EventLogs {
		log.OnProcessMessageComplete(message, reply)
	}
}

// OnProcessMessage implements middleware.EventLog.
func (l MultiEventLog) OnProcessMessage(message *redis.Message) {
	for _, log := range l.EventLogs {
		log.OnProcessMessage(message)
	}
}
