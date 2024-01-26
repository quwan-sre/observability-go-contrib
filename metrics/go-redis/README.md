## Redis client for Go

go-redis 对外提供了 Hook 机制，用于挂载 Middleware 扩展额外的请求响应逻辑。

```go
type Hook interface {
	DialHook(next DialHook) DialHook
	ProcessHook(next ProcessHook) ProcessHook
	ProcessPipelineHook(next ProcessPipelineHook) ProcessPipelineHook
}

type (
	DialHook			func(ctx context.Context, network, addr string) (net.Conn, error)
	ProcessHook		 func(ctx context.Context, cmd Cmder) error
	ProcessPipelineHook func(ctx context.Context, cmds []Cmder) error
)
```

Client 在 Init 时会初始化基础的 Hook，用户可以通过 AddHook 方法填加自定义的 Hook：
```go
func (c *Client) init() {
	c.cmdable = c.Process
	c.initHooks(hooks{
		dial:		c.baseClient.dial,
		process:	c.baseClient.process,
		pipeline:	c.baseClient.processPipeline,
		txPipeline:	c.baseClient.processTxPipeline,
	})
}

func (hs *hooksMixin) AddHook(hook Hook) {
	hs.slice = append(hs.slice, hook)
	hs.chain()
}
```

go-contrib 中实现了 `Hook` Interface，通过以下代码可以在 Client 初始化后添加：
```go
package main

import (
	"github.com/go-redis/redis/v9"

	metrics "github.com/quwan-sre/observability-go-contrib/metrics/go-redis"
)

func main() {

	redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{"127.0.0.1:6379"},
	})

	// 填加 Hook
	redisClient.AddHook(metrics.NewRedisHook([]string{"127.0.0.1:6379"}))
	return
}
```

由于 Hook 方法执行时想要获取 Redis Server 的信息困难，因此在 `NewRedisHook` 方法中需要提供连接 IP:Port，这些信息将直接用于填充指标的 Label。 