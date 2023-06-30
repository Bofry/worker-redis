package test

import (
	"context"
	"fmt"
	"log"

	"github.com/Bofry/host"
	"github.com/Bofry/trace"
	"go.opentelemetry.io/otel/propagation"
)

var (
	_ host.App                    = new(App)
	_ host.AppStaterConfigurator  = new(App)
	_ host.AppTracingConfigurator = new(App)
)

type App struct {
	Host            *Host
	Config          *Config
	ServiceProvider *ServiceProvider

	Component       *MockComponent
	ComponentRunner *MockComponentRunner
}

func (app *App) Init() {
	fmt.Println("App.Init()")

	app.Component = &MockComponent{}
	app.ComponentRunner = &MockComponentRunner{prefix: "MockComponentRunner"}
}

func (app *App) OnInit() {
}

func (app *App) OnInitComplete() {
}

func (app *App) OnStart(ctx context.Context) {
}

func (app *App) OnStop(ctx context.Context) {
	{
		defaultLogger.Printf("stoping TracerProvider")
		tp := trace.GetTracerProvider()
		err := tp.Shutdown(ctx)
		if err != nil {
			defaultLogger.Printf("stoping TracerProvider error: %+v", err)
		}
	}
}

func (app *App) ConfigureLogger(l *log.Logger) {
	l.SetFlags(defaultLogger.Flags())
	l.SetOutput(defaultLogger.Writer())
}

func (app *App) Logger() *log.Logger {
	return defaultLogger
}

func (app *App) ConfigureTracerProvider() {
	if len(app.Config.JaegerTraceUrl) == 0 {
		tp, _ := trace.NoopProvider()
		trace.SetTracerProvider(tp)
		return
	}

	tp, err := trace.JaegerProvider(app.Config.JaegerTraceUrl,
		trace.ServiceName("redis-trace-demo"),
		trace.Environment("go-bofry-worker-redis-test"),
		trace.Pid(),
	)
	if err != nil {
		defaultLogger.Fatal(err)
	}

	trace.SetTracerProvider(tp)
}

func (app *App) TracerProvider() *trace.SeverityTracerProvider {
	return trace.GetTracerProvider()
}

func (app *App) ConfigureTextMapPropagator() {
	trace.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func (app *App) TextMapPropagator() propagation.TextMapPropagator {
	return trace.GetTextMapPropagator()
}
