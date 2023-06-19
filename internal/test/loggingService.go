package test

import (
	"log"

	redis "github.com/Bofry/worker-redis"
)

var _ redis.LoggingService = new(LoggingService)

type LoggingService struct {
	logger *log.Logger
}

// ConfigureLogger implements middleware.LoggingService.
func (s *LoggingService) ConfigureLogger(l *log.Logger) {
	s.logger = l
}

// CreateEventLog implements middleware.LoggingService.
func (s *LoggingService) CreateEventLog(ev redis.EventEvidence) redis.EventLog {
	s.logger.Println("CreateEventLog()")
	return EventLog{
		logger:   s.logger,
		evidence: ev,
	}
}
