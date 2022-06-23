package standardtest

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Bofry/config"
	redis "github.com/Bofry/worker-redis"
)

func TestStarter(t *testing.T) {
	err := setupTestStarter()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := teardownTestStarter()
		if err != nil {
			t.Fatal(err)
		}
	}()

	app := App{}
	starter := redis.Startup(&app).
		Middlewares(
			redis.UseStreamGateway(&StreamGateway{}),
			redis.UseErrorHandler(func(err error) (disposed bool) {
				t.Logf("catch err: %v", err)
				return false
			}),
		).
		ConfigureConfiguration(func(service *config.ConfigurationService) {
			service.
				LoadEnvironmentVariables("").
				LoadYamlFile("config.yaml").
				LoadCommandArguments()

			t.Logf("%+v\n", app.Config)
		})

	runCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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

func setupTestStarter() error {
	opt := &redis.UniversalOptions{
		Addrs: []string{os.Getenv("REDIS_SERVER")},
		DB:    0,
	}

	admin, err := redis.NewAdminClient(opt)
	if err != nil {
		return err
	}
	defer admin.Close()

	{
		/*
			DEL gotestStream
		*/
		_, err = admin.Handle().Del("gotestStream").Result()
		if err != nil {
			return err
		}

		/*
			XGROUP CREATE gotestStream gotestGroup $ MKSTREAM

			XADD gotestStream * name luffy age 19
			XADD gotestStream * name nami age 21

		*/
		_, err = admin.CreateConsumerGroupAndStream("gotestStream", "gotestGroup", redis.StreamLastDeliveredID)
		if err != nil {
			return err
		}

		p, err := redis.NewForwarder(opt)
		if err != nil {
			return err
		}
		defer p.Close()

		_, err = p.Write("gotestStream", redis.StreamAsteriskID, map[string]interface{}{
			"name": "luffy",
			"age":  19,
		})
		if err != nil {
			return err
		}
		_, err = p.Write("gotestStream", redis.StreamAsteriskID, map[string]interface{}{
			"name": "nami",
			"age":  21,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func teardownTestStarter() error {
	admin, err := redis.NewAdminClient(&redis.UniversalOptions{
		Addrs: []string{os.Getenv("REDIS_SERVER")},
		DB:    0,
	})
	if err != nil {
		return err
	}
	defer admin.Close()

	{
		/*
			XGROUP DESTROY gotestStream gotestGroup
		*/
		_, err = admin.DeleteConsumerGroup("gotestStream", "gotestGroup")
		if err != nil {
			return err
		}

		/*
			DEL gotestStream
		*/
		_, err = admin.Handle().Del("gotestStream").Result()
		if err != nil {
			return err
		}
	}
	return nil
}
