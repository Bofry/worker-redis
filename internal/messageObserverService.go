package internal

import (
	"fmt"
	"reflect"
	"sync"
)

type MessageObserverService struct {
	MessageObservers map[reflect.Type]MessageObserver

	handlerObservers            map[string][]MessageObserver
	handlerObserversInitializer sync.Once
}

func (s *MessageObserverService) RegisterMessageObservers(msg *Message, handlerID string) {
	if s.handlerObservers != nil {
		observers := s.handlerObservers[handlerID]
		if len(observers) > 0 {
			d, ok := msg.Delegate.(*ContextMessageDelegate)
			if ok {
				d.registerMessageObservers(observers)
			}
		}
	}
}

func (s *MessageObserverService) UnregisterAllMessageObservers(msg *Message) {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if ok {
		d.unregisterAllMessageObservers()
	}
}

func (s *MessageObserverService) init(messageManager interface{}) {
	if len(s.MessageObservers) == 0 {
		return
	}
	if messageManager == nil {
		return
	}

	s.initHandlerObserverMap()
	s.buildHandlerObservers(messageManager)
}

func (s *MessageObserverService) initHandlerObserverMap() {
	s.handlerObserversInitializer.Do(func() {
		s.handlerObservers = make(map[string][]MessageObserver)
	})
}

func (s *MessageObserverService) buildHandlerObservers(messageManager interface{}) {
	var (
		rvManager reflect.Value = reflect.ValueOf(messageManager)
	)
	if rvManager.Kind() != reflect.Pointer || rvManager.IsNil() {
		return
	}

	rvManager = reflect.Indirect(rvManager)
	numOfHandles := rvManager.NumField()
	for i := 0; i < numOfHandles; i++ {
		rvHandler := rvManager.Field(i)
		if rvHandler.Kind() != reflect.Pointer || rvHandler.IsNil() {
			continue
		}

		var observers []MessageObserver

		if isMessageObserverAffinity(rvHandler) {
			affinity := asMessageObserverAffinity(rvHandler)
			if affinity != nil {
				types := affinity.MessageObserverTypes()
				for _, t := range types {
					v := s.MessageObservers[t]
					if v == nil {
						panic(fmt.Sprintf("cannot find MessageObserver of type %s", t.String()))
					}
					observers = append(observers, v)
				}
			}
		}

		rvHandler = reflect.Indirect(rvHandler)
		if rvHandler.Kind() == reflect.Struct {
			info := rvManager.Type().Field(i)
			if _, ok := s.handlerObservers[info.Name]; !ok {
				s.registerMessageObserver(info.Name, observers)
			}
		}
	}
}

func (s *MessageObserverService) registerMessageObserver(id string, observers []MessageObserver) {
	container := s.handlerObservers

	if observers != nil {
		if _, ok := container[id]; ok {
			RedisWorkerLogger.Fatalf("specified id '%s' already exists", id)
		}
		container[id] = observers
	}
}
