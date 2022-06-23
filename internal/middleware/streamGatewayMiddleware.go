package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/worker-redis/internal"
)

var _ host.Middleware = new(StreamGatewayMiddleware)

type StreamGatewayMiddleware struct {
	StreamGateway interface{}
}

func (m *StreamGatewayMiddleware) Init(appCtx *host.AppContext) {
	var (
		kafkaworker = asRedisWorker(appCtx.Host())
		preparer    = internal.NewRedisWorkerPreparer(kafkaworker)
	)

	binder := &StreamGatewayBinder{
		router:                           preparer.Router(),
		appContext:                       appCtx,
		configureUnhandledMessageHandler: preparer.RegisterUnhandledMessageHandler,
		configureStream:                  preparer.RegisterStream,
	}

	err := m.performBindTopicGateway(m.StreamGateway, binder)
	if err != nil {
		panic(err)
	}
}

func (m *StreamGatewayMiddleware) performBindTopicGateway(target interface{}, binder *StreamGatewayBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagName:     TAG_STREAM,
			TagResolver: StreamTagResolve,
		},
	)
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}
