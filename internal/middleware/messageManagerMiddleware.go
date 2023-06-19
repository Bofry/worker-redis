package middleware

import (
	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	. "github.com/Bofry/worker-redis/internal"
)

var _ host.Middleware = new(MessageManagerMiddleware)

type MessageManagerMiddleware struct {
	MessageManager interface{}
}

func (m *MessageManagerMiddleware) Init(app *host.AppModule) {
	var (
		kafkaworker = asRedisWorker(app.Host())
		registrar   = NewRedisWorkerRegistrar(kafkaworker)
	)

	// register RequestManager offer FasthttpHost processing later.
	registrar.SetMessageManager(m.MessageManager)

	// binding MessageManage
	binder := &MessageManagerBinder{
		registrar: registrar,
		app:       app,
	}

	err := m.performBindTopicGateway(m.MessageManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *MessageManagerMiddleware) performBindTopicGateway(target interface{}, binder *MessageManagerBinder) error {
	prototype, err := structproto.Prototypify(target,
		&structproto.StructProtoResolveOption{
			TagName:             TAG_STREAM,
			TagResolver:         StreamTagResolve,
			CheckDuplicateNames: true,
		},
	)
	if err != nil {
		return err
	}

	return prototype.Bind(binder)
}
