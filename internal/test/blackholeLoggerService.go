package test

import (
	"bytes"
	"log"

	redis "github.com/Bofry/worker-redis"
)

var _ redis.LoggingService = new(BlackholeLoggerService)

type BlackholeLoggerService struct {
	Buffer *bytes.Buffer
}

func (s *BlackholeLoggerService) CreateEventLog(ev redis.EventEvidence) redis.EventLog {
	s.Buffer.WriteString("CreateEventLog()")
	s.Buffer.WriteByte('\n')
	return &BlackholeEventLog{
		buffer: s.Buffer,
	}
}

func (*BlackholeLoggerService) ConfigureLogger(l *log.Logger) {
}

var _ redis.EventLog = new(BlackholeEventLog)

type BlackholeEventLog struct {
	buffer *bytes.Buffer
}

func (l *BlackholeEventLog) OnError(message *redis.Message, err interface{}, stackTrace []byte) {
	l.buffer.WriteString("LogError()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) OnProcessMessageComplete(message *redis.Message, reply redis.ReplyCode) {
	l.buffer.WriteString("OnProcessMessageComplete()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) OnProcessMessage(message *redis.Message) {
	l.buffer.WriteString("OnProcessMessage()")
	l.buffer.WriteByte('\n')
}

func (l *BlackholeEventLog) Flush() {
	l.buffer.WriteString("Flush()")
	l.buffer.WriteByte('\n')
}
