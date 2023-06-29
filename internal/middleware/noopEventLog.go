package middleware

import (
	redis "github.com/Bofry/lib-redis-stream"
	"github.com/Bofry/worker-redis/internal"
)

var _ EventLog = NoopEventLog(0)

type NoopEventLog int

// Flush implements EventLog.
func (NoopEventLog) Flush() {}

// OnError implements EventLog.
func (NoopEventLog) OnError(*redis.Message, interface{}, []byte) {}

// OnProcessMessage implements EventLog.
func (NoopEventLog) OnProcessMessage(*redis.Message) {}

// OnProcessMessageComplete implements EventLog.
func (NoopEventLog) OnProcessMessageComplete(*redis.Message, internal.ReplyCode) {}
