## Gin Web Framework

Gin 对外提供了 Middleware 机制，用于挂载 Middleware 扩展额外的请求响应逻辑。

```go
// Use attaches a global middleware to the router. i.e. the middleware attached through Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}
```

在 Gin 框架中可以通过以下代码使用 go-contrib 埋点：
```go
package main

import (
	"github.com/gin-gonic/gin"

	metrics "github.com/quwan-sre/observability-go-contrib/metrics/gin"
)

func main() {
	r := gin.New()

	// 添加 Middleware
	r.Use(metrics.NewMetricsMiddleware())

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
		return
	})
}
```
