package internal

import redis "github.com/Bofry/lib-redis-stream"

var _ redis.MessageDelegate = NoopMessageDelegate(0)

type NoopMessageDelegate int

// OnAck implements redis.MessageDelegate.
func (n NoopMessageDelegate) OnAck(msg *redis.Message) {}

// OnDel implements redis.MessageDelegate.
func (n NoopMessageDelegate) OnDel(msg *redis.Message) {}
