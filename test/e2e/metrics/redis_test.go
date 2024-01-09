package metrics

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	go_redis "github.com/quwan-sre/observability-go-contrib/metrics/go-redis"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	redisClient redis.UniversalClient
	redisHost   = "127.0.0.1"
	redisPort   = "6379"
)

func initRedisClient() {
	fmt.Println("Initializing redis client...")
	host, port := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
	if host != "" {
		redisHost = host
	}
	if port != "" {
		redisPort = port
	}

	redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{strings.Join([]string{redisHost, redisPort}, ":")},
	})
	redisClient.AddHook(go_redis.NewRedisHook([]string{strings.Join([]string{redisHost, redisPort}, ":")}))
	return
}

func TestRedisGo(t *testing.T) {
	initRedisClient()

	health := false
	for i := 0; i < 3; i++ {
		if err := redisClient.Ping(context.TODO()).Err(); err != nil {
			fmt.Printf("redis health check: %v\n", err)
			time.Sleep(time.Second)
		} else {
			health = true
			break
		}
	}

	if !health {
		t.Fatalf("redis not ready")
	}

	if redisClient.Set(context.TODO(), "foo", "bar", 0).Err() != nil {
		t.Error("redis err")
	}
	redisClient.Get(context.TODO(), "hello")
	redisClient.Do(context.TODO(), "not_exist_command", "invalid_value")

	pipe := redisClient.Pipeline()
	pipe.Del(context.TODO(), "foo")
	pipe.Get(context.TODO(), "foo")
	pipe.Del(context.TODO(), "foo")
	pipe.Do(context.TODO(), "not_exist_command", "invalid_value")
	pipe.Exec(context.TODO())

	return
}
