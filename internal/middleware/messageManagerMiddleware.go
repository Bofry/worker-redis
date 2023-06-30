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

// Init implements internal.Middleware.
func (m *MessageManagerMiddleware) Init(app *host.AppModule) {
	var (
		worker    = asRedisWorker(app.Host())
		registrar = NewRedisWorkerRegistrar(worker)
	)

	// register MessageManager offer RedisWorker processing later.
	registrar.SetMessageManager(m.MessageManager)

	// binding MessageManager
	binder := &MessageManagerBinder{
		registrar: registrar,
		app:       app,
	}

	err := m.bindMessageManager(m.MessageManager, binder)
	if err != nil {
		panic(err)
	}
}

func (m *MessageManagerMiddleware) bindMessageManager(target interface{}, binder *MessageManagerBinder) error {
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
