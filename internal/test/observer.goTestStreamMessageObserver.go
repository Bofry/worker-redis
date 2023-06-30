package test

import (
	"fmt"
	"reflect"

	redis "github.com/Bofry/worker-redis"
	"github.com/Bofry/worker-redis/tracing"
)

var _ redis.MessageObserver = new(GoTestStreamMessageObserver)

type GoTestStreamMessageObserver struct {
	ServiceProvider *ServiceProvider
}

func (*GoTestStreamMessageObserver) Init() {
	fmt.Println("GoTestStreamMessageObserver.Init()")
}

// OnAck implements internal.MessageObserver.
func (o *GoTestStreamMessageObserver) OnAck(ctx *redis.Context, message *redis.Message) {
	tr := tracing.GetTracer(o)
	sp := tr.Start(ctx, "OnAck()")
	defer sp.End()

	o.ServiceProvider.Logger().Println("GoTestStreamMessageObserver.OnAck()")
}

// OnDel implements internal.MessageObserver.
func (o *GoTestStreamMessageObserver) OnDel(ctx *redis.Context, message *redis.Message) {
	o.ServiceProvider.Logger().Println("GoTestStreamMessageObserver.OnDel()")
}

func (o *GoTestStreamMessageObserver) Type() reflect.Type {
	return reflect.TypeOf(o)
}
