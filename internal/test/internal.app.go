package test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Bofry/host"
	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
	"go.opentelemetry.io/otel/propagation"
)

var (
	defaultLogger *log.Logger = log.New(log.Writer(), "[worker-redis-test] ", log.LstdFlags|log.Lmsgprefix|log.LUTC)
)

var (
	_ host.App                    = new(App)
	_ host.AppStaterConfigurator  = new(App)
	_ host.AppTracingConfigurator = new(App)
)

type (
	App struct {
		Host            *Host
		Config          *Config
		ServiceProvider *ServiceProvider

		Component       *MockComponent
		ComponentRunner *MockComponentRunner
	}

	Host redis.Worker

	Config struct {
		// redis
		RedisAddresses           []string      `env:"*REDIS_SERVER"        yaml:"-"`
		RedisConsumerGroup       string        `env:"-"                    yaml:"RedisConsumerGroup"`
		RedisConsumerName        string        `env:"-"                    yaml:"RedisConsumerName"`
		RedisMaxInFlight         int64         `env:"-"                    yaml:"RedisMaxInFlight"`
		RedisMaxPollingTimeout   time.Duration `env:"-"                    yaml:"RedisMaxPollingTimeout"`
		RedisClaimMinIdleTime    time.Duration `env:"-"                    yaml:"RedisClaimMinIdleTime"`
		RedisIdlingTimeout       time.Duration `env:"-"                    yaml:"RedisIdlingTimeout"`
		RedisClaimSensitivity    int           `env:"-"                    yaml:"RedisClaimSensitivity"`
		RedisClaimOccurrenceRate int32         `env:"-"                    yaml:"RedisClaimOccurrenceRate"`

		// jaeger
		JaegerTraceUrl string `yaml:"jaegerTraceUrl"`
		JaegerQueryUrl string `yaml:"jaegerQueryUrl"`
	}

	ServiceProvider struct {
		ResourceName string
	}
)

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

func (provider *ServiceProvider) Init(conf *Config) {
	fmt.Println("ServiceProvider.Init()")
	provider.ResourceName = "demo resource"
}

func (h *Host) Init(conf *Config) {
	h.RedisOption = &redis.UniversalOptions{
		Addrs: conf.RedisAddresses,
	}
	h.ConsumerGroup = conf.RedisConsumerGroup
	h.ConsumerName = conf.RedisConsumerName
	h.MaxInFlight = conf.RedisMaxInFlight
	h.MaxPollingTimeout = conf.RedisMaxPollingTimeout
	h.ClaimMinIdleTime = conf.RedisClaimMinIdleTime
	h.IdlingTimeout = conf.RedisIdlingTimeout
	h.ClaimSensitivity = conf.RedisClaimSensitivity
	h.ClaimOccurrenceRate = conf.RedisClaimOccurrenceRate
	h.AllowCreateGroup = true
}

func (h *Host) OnError(err error) (disposed bool) {
	return false
}
