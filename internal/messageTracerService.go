package internal

import (
	"reflect"
	"sync"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type MessageTracerService struct {
	TracerManager *TracerManager

	Enabled bool

	InvalidMessageHandlerComponentID string

	invalidMessageTracer *trace.SeverityTracer

	tracerProvider    *trace.SeverityTracerProvider
	textMapPropagator propagation.TextMapPropagator

	tracers            map[string]*trace.SeverityTracer
	tracersInitializer sync.Once
}

func NewMessageTracerService() *MessageTracerService {
	return &MessageTracerService{}
}

func (s *MessageTracerService) Tracer(id string) *trace.SeverityTracer {
	if s.tracers != nil {
		if tr, ok := s.tracers[id]; ok {
			return tr
		}
	}
	return s.invalidMessageTracer
}

func (s *MessageTracerService) TextMapPropagator() propagation.TextMapPropagator {
	return s.TracerManager.TextMapPropagator
}

func (s *MessageTracerService) init(messageManager interface{}) {
	if s.Enabled {
		if s.tracerProvider != nil {
			s.TracerManager.TracerProvider = s.tracerProvider
		}
		if s.textMapPropagator != nil {
			s.TracerManager.TextMapPropagator = s.textMapPropagator
		}

		s.makeTracerMap()
		s.buildTracer(messageManager)
	}
	s.makeInvalidMessageTracer()
}

func (s *MessageTracerService) makeTracerMap() {
	s.tracersInitializer.Do(func() {
		s.tracers = make(map[string]*trace.SeverityTracer)
	})
}

func (s *MessageTracerService) buildTracer(requestManager interface{}) {
	var (
		rvManager reflect.Value = reflect.ValueOf(requestManager)
	)
	if rvManager.Kind() != reflect.Pointer || rvManager.IsNil() {
		return
	}

	rvManager = reflect.Indirect(rvManager)
	numOfHandles := rvManager.NumField()
	for i := 0; i < numOfHandles; i++ {
		rvRequest := rvManager.Field(i)
		if rvRequest.Kind() != reflect.Pointer || rvRequest.IsNil() {
			continue
		}

		rvRequest = reflect.Indirect(rvRequest)
		if rvRequest.Kind() == reflect.Struct {
			tracer := s.TracerManager.createManagedTracer(rvRequest.Type())

			info := rvManager.Type().Field(i)
			if _, ok := s.tracers[info.Name]; !ok {
				s.registerTracer(info.Name, tracer)
			}
		}
	}
}

func (s *MessageTracerService) registerTracer(id string, tracer *trace.SeverityTracer) {
	container := s.tracers

	if tracer != nil {
		if _, ok := container[id]; ok {
			RedisWorkerLogger.Fatalf("specified id '%s' already exists", id)
		}
		container[id] = tracer
	}
}

func (s *MessageTracerService) makeInvalidMessageTracer() {
	if len(s.InvalidMessageHandlerComponentID) > 0 {
		v, ok := s.tracers[s.InvalidMessageHandlerComponentID]
		if ok {
			s.invalidMessageTracer = v
			return
		}
	}
	s.invalidMessageTracer = s.TracerManager.createTracer("")
}
