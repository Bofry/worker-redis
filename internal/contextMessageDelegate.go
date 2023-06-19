package internal

import (
	"sync"

	redis "github.com/Bofry/lib-redis-stream"
)

var _ redis.MessageDelegate = new(ContextMessageDelegate)

type ContextMessageDelegate struct {
	parent redis.MessageDelegate

	ctx *Context

	mu sync.Mutex
}

func NewContextMessageDelegate(ctx *Context) *ContextMessageDelegate {
	return &ContextMessageDelegate{
		ctx: ctx,
	}
}

// OnAck implements redis.MessageDelegate.
func (d *ContextMessageDelegate) OnAck(msg *redis.Message) {
	d.parent.OnAck(msg)
	GlobalContextHelper.InjectReplyCode(d.ctx, PASS)
}

// OnDel implements redis.MessageDelegate.
func (d *ContextMessageDelegate) OnDel(msg *redis.Message) {
	d.parent.OnDel(msg)
}

func (d *ContextMessageDelegate) configure(msg *redis.Message) {
	if d.parent == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if d.parent == nil {
			d.parent = msg.Delegate
			msg.Delegate = d
		}
	}
}
