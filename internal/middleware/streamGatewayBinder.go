package middleware

import (
	"fmt"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/tagresolver"
	"github.com/Bofry/structproto/util/reflectutil"
	"github.com/Bofry/worker-redis/internal"
)

var _ structproto.StructBinder = new(StreamGatewayBinder)

type StreamGatewayBinder struct {
	appContext                       *host.AppContext
	router                           internal.Router
	configureUnhandledMessageHandler ConfigureUnhandledMessageHandleProc
	configureStream                  ConfigureStream
}

func (b *StreamGatewayBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *StreamGatewayBinder) Bind(field structproto.FieldInfo, rv reflect.Value) error {
	if !rv.IsValid() {
		return fmt.Errorf("specifiec argument 'rv' is invalid")
	}

	// assign zero if rv is nil
	rvMessageHandler := reflectutil.AssignZero(rv)
	binder := &MessageHandlerBinder{
		messageHandlerType: rv.Type().Name(),
		components: map[string]reflect.Value{
			host.APP_CONFIG_FIELD:           b.appContext.Config(),
			host.APP_SERVICE_PROVIDER_FIELD: b.appContext.ServiceProvider(),
		},
	}
	err := b.preformBindMessageHandler(rvMessageHandler, binder)
	if err != nil {
		return err
	}

	// register MessageHandlers
	return b.registerRoute(field.Name(), rvMessageHandler)
}

func (b *StreamGatewayBinder) Deinit(context *structproto.StructProtoContext) error {
	return nil
}

func (b *StreamGatewayBinder) preformBindMessageHandler(target reflect.Value, binder *MessageHandlerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagResolver: tagresolver.NoneTagResolver,
		})
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}

func (b *StreamGatewayBinder) registerRoute(stream string, rv reflect.Value) error {
	// register MessageHandlers
	if isMessageHandler(rv) {
		handler := asMessageHandler(rv)
		if handler != nil {
			if stream == UNHANDLED_MESSAGE_HANDLER_TOPIC_SYMBOL {
				b.configureUnhandledMessageHandler(handler)
			} else {
				// TODO....
				b.configureStream(internal.StreamOffset{
					Stream: stream,
				})
				b.router.Add(stream, handler)
			}
		}
	}
	return nil
}
