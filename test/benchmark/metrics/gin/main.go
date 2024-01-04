package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/quwan-sre/observability-go-contrib/metrics/gin"
)

const (
	httpPort           = 8081
	httpInstrumentPort = 8082
)

func main() {
	go newHTTPServer()
	go newInstrumentedHTTPServer()
}

func newHTTPServer() {
	r := gin.New()

	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
		return
	})

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
		return
	})

	r.Run(fmt.Sprintf("127.0.0.1:%d", httpPort))
}

func newInstrumentedHTTPServer() {
	r := gin.New()

	r.Use(metrics.NewMetricsMiddleware())

	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
		return
	})

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
		return
	})

	r.Run(fmt.Sprintf("127.0.0.1:%d", httpInstrumentPort))
}
