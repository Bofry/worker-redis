package test

import (
	"fmt"
	"log"

	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

type ServiceProvider struct {
	ResourceName string

	TestStreamMessageCounter *TestStreamMessageCounter
}

func (p *ServiceProvider) Init(conf *Config) {
	fmt.Println("ServiceProvider.Init()")
	p.ResourceName = "demo resource"

	p.TestStreamMessageCounter = new(TestStreamMessageCounter)
}

func (p *ServiceProvider) TracerProvider() *trace.SeverityTracerProvider {
	return trace.GetTracerProvider()
}

func (p *ServiceProvider) TextMapPropagator() propagation.TextMapPropagator {
	return trace.GetTextMapPropagator()
}

func (p *ServiceProvider) Logger() *log.Logger {
	return defaultLogger
}
