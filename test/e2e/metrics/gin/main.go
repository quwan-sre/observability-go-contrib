package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/quwan-sre/observability-go-contrib/metrics/gin"
)

func main() {
	r := gin.Default()

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

	r.GET("/exist", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
		return
	})

	r.Run()
}
