package middleware

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/reflecting"
	"github.com/Bofry/structproto/tagresolver"
	"github.com/Bofry/worker-redis/internal"
)

var _ structproto.StructBinder = new(MessageManagerBinder)

type MessageManagerBinder struct {
	registrar *internal.RedisWorkerRegistrar
	app       *host.AppModule
}

func (b *MessageManagerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *MessageManagerBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvMessageHandler := reflecting.AssignZero(rv)
	binder := &MessageHandlerBinder{
		messageHandlerType: rv.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.app.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.app.ServiceProvider(),
		},
	}
	err := b.preformBindMessageHandler(rvMessageHandler, binder)
	if err != nil {
		return err
	}

	// register MessageHandlers
	var (
		moduleID = field.IDName()
		stream   = field.Name()
		offset   = field.Tag().Get(TAG_OFFSET)
	)

	if !b.isKnownStream(stream) {
		optExpandEnv := field.Tag().Get(TAG_OPT_EXPAND_ENV)
		if optExpandEnv != "off" || len(optExpandEnv) == 0 || optExpandEnv == "on" {
			stream = os.ExpandEnv(stream)
		}
	}

	return b.registerRoute(moduleID, stream, offset, rvMessageHandler)
}

func (b *MessageManagerBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *MessageManagerBinder) preformBindMessageHandler(target reflect.Value, binder *MessageHandlerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *MessageManagerBinder) registerRoute(moduleID, stream, offset string, rv reflect.Value) error {
	// register MessageHandlers
	if isMessageHandler(rv) {
		handler := asMessageHandler(rv)
		if handler != nil {
			if stream == INVALID_MESSAGE_HANDLER_TOPIC_SYMBOL {
				b.registrar.SetInvalidMessageHandler(handler)
			} else {
				if offset != "-" {
					b.registrar.RegisterStream(internal.StreamOffset{
						Stream: stream,
						Offset: offset,
					})
				}
				b.registrar.AddRouter(stream, handler, moduleID)
			}
		}
	}
	return nil
}

func (b *MessageManagerBinder) isKnownStream(stream string) bool {
	switch stream {
	case INVALID_MESSAGE_HANDLER_TOPIC_SYMBOL:
		return true
	}
	return false
}
