package middleware

import (
	"fmt"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/reflecting"
	"github.com/Bofry/structproto/tagresolver"
	"github.com/Bofry/worker-redis/internal"
)

var _ structproto.StructBinder = new(MessageObserverManagerBinder)

type MessageObserverManagerBinder struct {
	registrar *internal.RedisWorkerRegistrar
	app       *host.AppModule
}

// Init implements structproto.StructBinder.
func (*MessageObserverManagerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

// Bind implements structproto.StructBinder.
func (b *MessageObserverManagerBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvMessageObserver := reflecting.AssignZero(rv)
	binder := &MessageObserverBinder{
		messageHandlerType: rv.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.app.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.app.ServiceProvider(),
		},
	}
	err := b.bindMessageObserver(rvMessageObserver, binder)
	if err != nil {
		return err
	}

	// register MessageHandlers
	var (
		moduleID = field.IDName()
		stream   = field.Name()
		offset   = field.Tag().Get(TAG_OFFSET)
	)

	return b.registerObservers(moduleID, stream, offset, rvMessageObserver)
}

// Deinit implements structproto.StructBinder.
func (*MessageObserverManagerBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *MessageObserverManagerBinder) bindMessageObserver(target reflect.Value, binder *MessageObserverBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *MessageObserverManagerBinder) registerObservers(moduleID, stream, offset string, rv reflect.Value) error {
	// register MessageObservers
	if isMessageObserver(rv) {
		observer := asMessageObserver(rv)
		if observer != nil {
			b.registrar.RegisterMessageObserver(observer)
		}
	}
	return nil
}
