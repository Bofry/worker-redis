package internal

import (
	"sync"
	"sync/atomic"

	redis "github.com/Bofry/lib-redis-stream"
)

var _ redis.MessageDelegate = new(ContextMessageDelegate)

type ContextMessageDelegate struct {
	parent redis.MessageDelegate

	ctx *Context

	messageObserver MessageObserver

	restricted int32
	mu         sync.Mutex
}

func NewContextMessageDelegate(ctx *Context) *ContextMessageDelegate {
	return &ContextMessageDelegate{
		ctx: ctx,
	}
}

// OnAck implements redis.MessageDelegate.
func (d *ContextMessageDelegate) OnAck(msg *redis.Message) {
	if d.isRestricted() {
		GlobalNoopMessageDelegate.OnAck(nil)
		return
	}

	d.parent.OnAck(msg)
	GlobalContextHelper.InjectReplyCodeSafe(d.ctx, PASS)

	// observer
	if d.messageObserver != nil {
		d.messageObserver.OnAck(d.ctx, msg)
	}
}

// OnDel implements redis.MessageDelegate.
func (d *ContextMessageDelegate) OnDel(msg *redis.Message) {
	if d.isRestricted() {
		GlobalNoopMessageDelegate.OnDel(nil)
		return
	}

	d.parent.OnDel(msg)

	// observer
	if d.messageObserver != nil {
		d.messageObserver.OnDel(d.ctx, msg)
	}
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

func (d *ContextMessageDelegate) isRestricted() bool {
	return atomic.LoadInt32(&d.restricted) == 1
}

func (d *ContextMessageDelegate) restrict() {
	atomic.StoreInt32(&d.restricted, 1)
}

func (d *ContextMessageDelegate) unrestrict() {
	atomic.StoreInt32(&d.restricted, 0)
}

func (d *ContextMessageDelegate) registerMessageObservers(observers []MessageObserver) {
	d.messageObserver = CompositeMessageObserver(observers)
}

func (d *ContextMessageDelegate) unregisterAllMessageObservers() {
	d.messageObserver = nil
}
