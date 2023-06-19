package internal

import "context"

type MessageHandleService struct {
	modules []MessageHandleModule
}

func NewMessageHandleService() *MessageHandleService {
	return &MessageHandleService{}
}

func (s *MessageHandleService) Register(module MessageHandleModule) {
	size := len(s.modules)
	if size > 0 {
		last := s.modules[size-1]

		// ignore all new successor if the last RequestHandleModule cannot accept successor
		if !last.CanSetSuccessor() {
			return
		}

		last.SetSuccessor(module)
	}
	s.modules = append(s.modules, module)
}

func (s *MessageHandleService) ProcessMessage(ctx *Context, message *Message, state ProcessingState, recover *Recover) {
	if handler := s.first(); handler != nil {
		handler.ProcessMessage(ctx, message, state, recover)
	}
}

func (s *MessageHandleService) first() MessageHandleModule {
	if len(s.modules) > 0 {
		return s.modules[0]
	}
	return nil
}

func (s *MessageHandleService) triggerStart(ctx context.Context) error {
	var err error

	for _, m := range s.modules {
		if err = m.OnStart(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *MessageHandleService) triggerStop(ctx context.Context) <-chan error {
	ch := make(chan error)

	go func() {
		defer close(ch)
		for _, m := range s.modules {
			RedisWorkerLogger.Printf("stopping %T", m)

			err := m.OnStop(ctx)
			if err != nil {
				ch <- &StopError{
					v:   m,
					err: err,
				}
			}
		}
	}()

	return ch
}
