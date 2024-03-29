package internal

import (
	redis "github.com/Bofry/lib-redis-stream"
)

var _ redis.MessageDelegate = RestrictedMessageDelegate(0)

type RestrictedMessageDelegate int

// OnAck implements redis.MessageDelegate.
func (RestrictedMessageDelegate) OnAck(*redis.Message) {
	panic(RestrictedOperationError("Message.Ack() cannot be called by restricted area"))
}

// OnDel implements redis.MessageDelegate.
func (RestrictedMessageDelegate) OnDel(*redis.Message) {
	panic(RestrictedOperationError("Message.Del() cannot be called by restricted area"))
}
