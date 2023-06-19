package test

import (
	"fmt"
	"log"

	redis "github.com/Bofry/worker-redis"
)

var _ redis.EventLog = EventLog{}

type EventLog struct {
	logger   *log.Logger
	evidence redis.EventEvidence
}

// AfterProcessMessage implements middleware.EventLog.
func (l EventLog) OnProcessMessageComplete(message *redis.Message, reply redis.ReplyCode) {
	traceID := fmt.Sprintf("%s-%s",
		l.evidence.ProcessingSpanID(),
		l.evidence.ProcessingSpanID())

	l.logger.Printf("EventLog.OnProcessMessageComplete(): (%s) %s\n", traceID, message.ID)
}

// BeforeProcessMessage implements middleware.EventLog.
func (l EventLog) OnProcessMessage(message *redis.Message) {
	traceID := fmt.Sprintf("%s-%s",
		l.evidence.ProcessingSpanID(),
		l.evidence.ProcessingSpanID())

	l.logger.Printf("EventLog.OnProcessMessage(): (%s) %s\n", traceID, message.ID)
}

// Flush implements middleware.EventLog.
func (l EventLog) Flush() {
	l.logger.Println("EventLog.Flush()")
}

// LogError implements middleware.EventLog.
func (l EventLog) OnError(message *redis.Message, err interface{}, stackTrace []byte) {
	l.logger.Printf("EventLog.OnError(): %v\n", err)
}
