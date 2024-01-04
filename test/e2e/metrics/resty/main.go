package resty

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"time"
)

func main() {
	r := gin.Default()

	r.GET("/metrics", func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
		return
	})

	r.GET("/health", func(ctx *gin.Context) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
		return
	})

	r.Run()
}
