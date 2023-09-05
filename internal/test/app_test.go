package test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Bofry/config"
	"github.com/Bofry/trace"
	redis "github.com/Bofry/worker-redis"
	goredis "github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
)

var (
	__TEST_REDIS_SERVER     string
	__TEST_JAEGER_TRACE_URL string

	__ENV_FILE        = ".env"
	__ENV_FILE_SAMPLE = ".env.sample"

	__CONFIG_YAML_FILE        = "config.yaml"
	__CONFIG_YAML_FILE_SAMPLE = "config.yaml.sample"
)

type MessageManager struct {
	GotestStream  *GoTestStreamMessageHandler `stream:"gotestStream"   offset:"$"   @ExpandEnv:"off"`
	GotestStream2 *GoTestStreamMessageHandler `stream:"gotestStream2"  offset:"$"`
	GotestStream3 *GoTestStreamMessageHandler `stream:"gotestStream3"  offset:"$"`
	GotestStream4 *GoTestStreamMessageHandler `stream:"gotestStream4"  offset:"$"   @MessageStateKeyPrefix:"mystate:"`
	Invalid       *InvalidMessageHandler      `stream:"?"`
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func execRedisCommand(client *goredis.Client, cmd string) *goredis.Cmd {
	var args []interface{}
	for _, argv := range strings.Split(cmd, " ") {
		args = append(args, argv)
	}
	return client.Do(args...)
}

func TestMain(m *testing.M) {
	var err error

	_, err = os.Stat(__CONFIG_YAML_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__CONFIG_YAML_FILE_SAMPLE, __CONFIG_YAML_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	_, err = os.Stat(__ENV_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			err = copyFile(__ENV_FILE_SAMPLE, __ENV_FILE)
			if err != nil {
				panic(err)
			}
		}
	}

	{
		f, err := os.Open(__ENV_FILE)
		if err != nil {
			panic(err)
		}
		env, err := godotenv.Parse(f)
		if err != nil {
			panic(err)
		}
		__TEST_REDIS_SERVER = env["TEST_REDIS_SERVER"]
		__TEST_JAEGER_TRACE_URL = env["TEST_JAEGER_TRACE_URL"]
	}
	{
		client := goredis.NewClient(&goredis.Options{
			Addr: __TEST_REDIS_SERVER,
			DB:   0,
		})
		if client == nil {
			panic("fail to create redis.Client")
		}
		defer client.Close()

		for _, cmd := range []string{
			"DEL gotestStream",
			"DEL gotestStream2",

			"XGROUP CREATE gotestStream gotestGroup $ MKSTREAM",
			"XGROUP CREATE gotestStream2 gotestGroup $ MKSTREAM",

			"XADD gotestStream * name luffy age 19",
			"XADD gotestStream * name nami age 21",
			"XADD gotestStream2 * name roger age ??",
			"XADD gotestStream2 * name ace age 22",
			"XADD gotestStream3 * name luffy age 19 header:foo bar",
			"XADD gotestStream3 * name nami age 21 header:foo bar",
			"XADD gotestStream4 * name roger age ?? mystate:foo bar",
			"XADD gotestStream4 * name ace age 22 mystate:foo bar",
		} {
			_, err := execRedisCommand(client, cmd).Result()
			if err != nil {
				panic(err)
			}
		}
		defer func() {
			for _, cmd := range []string{
				"XGROUP DESTROY gotestStream gotestGroup",
				"XGROUP DESTROY gotestStream2 gotestGroup",

				"DEL gotestStream",
				"DEL gotestStream2",
			} {
				_, err := execRedisCommand(client, cmd).Result()
				if err != nil {
					panic(err)
				}
			}
		}()

	}

	m.Run()
}

func TestStartup(t *testing.T) {
	app := App{}
	starter := redis.Startup(&app).
		Middlewares(
			redis.UseMessageManager(&MessageManager{}),
			redis.UseErrorHandler(func(ctx *redis.Context, message *redis.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}

	// assert app.Config
	{
		conf := app.Config
		var expectedRedisAddresses []string = []string{os.Getenv("REDIS_SERVER")}
		if !reflect.DeepEqual(conf.RedisAddresses, expectedRedisAddresses) {
			t.Errorf("assert 'Config.RedisAddress':: expected '%v', got '%v'", expectedRedisAddresses, conf.RedisAddresses)
		}
		var expectedRedisConsumerGroup string = "gotestGroup"
		if conf.RedisConsumerGroup != expectedRedisConsumerGroup {
			t.Errorf("assert 'Config.RedisConsumerGroup':: expected '%v', got '%v'", expectedRedisConsumerGroup, conf.RedisConsumerGroup)
		}
		var expectedRedisConsumerName string = "gotestConsumer"
		if conf.RedisConsumerName != expectedRedisConsumerName {
			t.Errorf("assert 'Config.RedisConsumerName':: expected '%v', got '%v'", expectedRedisConsumerName, conf.RedisConsumerName)
		}
		var expectedRedisMaxInFlight int64 = 8
		if conf.RedisMaxInFlight != expectedRedisMaxInFlight {
			t.Errorf("assert 'Config.RedisMaxInFlight':: expected '%v', got '%v'", expectedRedisMaxInFlight, conf.RedisMaxInFlight)
		}
		var expectedRedisMaxPollingTimeout time.Duration = 10 * time.Millisecond
		if conf.RedisMaxPollingTimeout != expectedRedisMaxPollingTimeout {
			t.Errorf("assert 'Config.RedisMaxPollingTimeout':: expected '%v', got '%v'", expectedRedisMaxPollingTimeout, conf.RedisMaxPollingTimeout)
		}
		var expectedRedisClaimMinIdleTime time.Duration = 30 * time.Second
		if conf.RedisClaimMinIdleTime != expectedRedisClaimMinIdleTime {
			t.Errorf("assert 'Config.RedisClaimMinIdleTime':: expected '%v', got '%v'", expectedRedisClaimMinIdleTime, conf.RedisClaimMinIdleTime)
		}
	}
}

func TestStartup_UseTracing(t *testing.T) {
	var (
		testStartAt time.Time

		traceIDStr string
		spanIDStr  string
	)

	{
		tp, err := trace.JaegerProvider(__TEST_JAEGER_TRACE_URL,
			trace.ServiceName("trace-demo"),
			trace.Environment("go-test"),
			trace.Pid(),
		)
		if err != nil {
			t.Fatal(err)
		}
		defer tp.Shutdown(context.Background())

		tr := tp.Tracer("worker-redis-test")
		sp := tr.Open(context.Background(), "TestStartup_UseTracing")
		defer sp.End()

		var (
			traceID = sp.TraceID()
			spanID  = sp.SpanID()
		)
		traceIDStr = hex.EncodeToString(traceID[:])
		spanIDStr = hex.EncodeToString(spanID[:])
	}

	client := goredis.NewClient(&goredis.Options{
		Addr: __TEST_REDIS_SERVER,
		DB:   0,
	})
	if client == nil {
		t.Fatal("fail to create redis.Client")
	}
	defer client.Close()
	cmd := fmt.Sprintf("XADD gotestStream4 * name zoro age 20 mystate:traceparent 00-%s-%s-01", traceIDStr, spanIDStr)
	_, err := execRedisCommand(client, cmd).Result()
	if err != nil {
		t.Fatal(err)
	}

	app := App{}
	starter := redis.Startup(&app).
		Middlewares(
			redis.UseMessageManager(&MessageManager{}),
			redis.UseErrorHandler(func(ctx *redis.Context, message *redis.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
			redis.UseTracing(true),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	testStartAt = time.Now()

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}

		testEndAt := time.Now()

		// wait 2 seconds
		time.Sleep(2 * time.Second)

		var queryUrl = fmt.Sprintf(
			"%s?end=%d&limit=50&lookback=1h&&service=redis-trace-demo&start=%d",
			app.Config.JaegerQueryUrl,
			testEndAt.UnixMicro(),
			testStartAt.UnixMicro())

		t.Log(queryUrl)
		req, err := http.NewRequest("GET", queryUrl, nil)
		if err != nil {
			t.Error(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert query 'Jeager Query Url StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		// t.Logf("%v", string(body))
		// parse content
		{
			var reply map[string]interface{}
			dec := json.NewDecoder(bytes.NewBuffer(body))
			dec.UseNumber()
			if err := dec.Decode(&reply); err != nil {
				t.Error(err)
			}

			data := reply["data"].([]interface{})
			if data == nil {
				t.Errorf("missing data section")
			}
			var expectedDataLength int = 9
			if expectedDataLength != len(data) {
				t.Errorf("assert 'Jaeger Query size of replies':: expected '%v', got '%v'", expectedDataLength, len(data))
			}
		}
	}
}

func TestStartup_UseLogging_And_UseTracing(t *testing.T) {
	var (
		testStartAt   time.Time
		loggingBuffer bytes.Buffer
	)

	app := App{}
	starter := redis.Startup(&app).
		Middlewares(
			redis.UseMessageManager(&MessageManager{}),
			redis.UseErrorHandler(func(ctx *redis.Context, message *redis.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
			redis.UseLogging(
				&LoggingService{},
				&BlackholeLoggerService{
					Buffer: &loggingBuffer,
				},
			),
			redis.UseTracing(true),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	testStartAt = time.Now()

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}

		testEndAt := time.Now()

		// wait 2 seconds
		time.Sleep(2 * time.Second)

		var queryUrl = fmt.Sprintf(
			"%s?end=%d&limit=50&lookback=1h&&service=redis-trace-demo&start=%d",
			app.Config.JaegerQueryUrl,
			testEndAt.UnixMicro(),
			testStartAt.UnixMicro())

		t.Log(queryUrl)
		req, err := http.NewRequest("GET", queryUrl, nil)
		if err != nil {
			t.Error(err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if resp.StatusCode != 200 {
			t.Errorf("assert query 'Jeager Query Url StatusCode':: expected '%v', got '%v'", 200, resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		// t.Logf("%v", string(body))
		// parse content
		{
			var reply map[string]interface{}
			dec := json.NewDecoder(bytes.NewBuffer(body))
			dec.UseNumber()
			if err := dec.Decode(&reply); err != nil {
				t.Error(err)
			}

			data := reply["data"].([]interface{})
			if data == nil {
				t.Errorf("missing data section")
			}
			var expectedDataLength int = 4
			if expectedDataLength != len(data) {
				t.Errorf("assert 'Jaeger Query size of replies':: expected '%v', got '%v'", expectedDataLength, len(data))
			}
		}

		// test loggingBuffer
		var expectedLoggingBuffer string = strings.Join([]string{
			"CreateEventLog()\n",
			"OnProcessMessage()\n",
			"OnProcessMessageComplete()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"OnProcessMessage()\n",
			"OnProcessMessageComplete()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"OnProcessMessage()\n",
			"OnProcessMessageComplete()\n",
			"Flush()\n",
			"CreateEventLog()\n",
			"OnProcessMessage()\n",
			"OnProcessMessageComplete()\n",
			"Flush()\n",
		}, "")
		if expectedLoggingBuffer != loggingBuffer.String() {
			t.Errorf("assert loggingBuffer:: expected '%v', got '%v'", expectedLoggingBuffer, loggingBuffer.String())
		}
	}
}

func TestStartup_UseMessageObserverManager(t *testing.T) {
	app := App{}
	starter := redis.Startup(&app).
		Middlewares(
			redis.UseMessageManager(&MessageManager{}),
			redis.UseErrorHandler(func(ctx *redis.Context, message *redis.Message, err interface{}) {
				t.Logf("catch err: %v", err)
			}),
			redis.UseMessageObserverManager(&MessageObserverManager),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := starter.Start(runCtx); err != nil {
		t.Error(err)
	}

	select {
	case <-runCtx.Done():
		if err := starter.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}

	// assert app.Config
	{
		conf := app.Config
		var expectedRedisAddresses []string = []string{os.Getenv("REDIS_SERVER")}
		if !reflect.DeepEqual(conf.RedisAddresses, expectedRedisAddresses) {
			t.Errorf("assert 'Config.RedisAddress':: expected '%v', got '%v'", expectedRedisAddresses, conf.RedisAddresses)
		}
		var expectedRedisConsumerGroup string = "gotestGroup"
		if conf.RedisConsumerGroup != expectedRedisConsumerGroup {
			t.Errorf("assert 'Config.RedisConsumerGroup':: expected '%v', got '%v'", expectedRedisConsumerGroup, conf.RedisConsumerGroup)
		}
		var expectedRedisConsumerName string = "gotestConsumer"
		if conf.RedisConsumerName != expectedRedisConsumerName {
			t.Errorf("assert 'Config.RedisConsumerName':: expected '%v', got '%v'", expectedRedisConsumerName, conf.RedisConsumerName)
		}
		var expectedRedisMaxInFlight int64 = 8
		if conf.RedisMaxInFlight != expectedRedisMaxInFlight {
			t.Errorf("assert 'Config.RedisMaxInFlight':: expected '%v', got '%v'", expectedRedisMaxInFlight, conf.RedisMaxInFlight)
		}
		var expectedRedisMaxPollingTimeout time.Duration = 10 * time.Millisecond
		if conf.RedisMaxPollingTimeout != expectedRedisMaxPollingTimeout {
			t.Errorf("assert 'Config.RedisMaxPollingTimeout':: expected '%v', got '%v'", expectedRedisMaxPollingTimeout, conf.RedisMaxPollingTimeout)
		}
		var expectedRedisClaimMinIdleTime time.Duration = 30 * time.Second
		if conf.RedisClaimMinIdleTime != expectedRedisClaimMinIdleTime {
			t.Errorf("assert 'Config.RedisClaimMinIdleTime':: expected '%v', got '%v'", expectedRedisClaimMinIdleTime, conf.RedisClaimMinIdleTime)
		}
	}
}
