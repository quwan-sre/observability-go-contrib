package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.ttyuyin.com/observability/go-contrib/metrics/common"
	"strconv"
	"time"
)

func NewMetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// before
		startTime := time.Now()

		// execute
		ctx.Next()

		// after
		latency := time.Since(startTime)
		endpoint := ""
		responseCode := -1

		if ctx != nil && ctx.Request != nil && ctx.Request.URL != nil {
			endpoint = ctx.Request.URL.Path
		}

		if ctx != nil && ctx.Writer != nil {
			responseCode = ctx.Writer.Status()
		}

		common.DefaultRPCReceiveRequestMetric.With(prometheus.Labels{
			"sdk":              common.SDKGin,
			"request_protocol": common.RequestProtocolHTTP,
			"endpoint":         endpoint,
			"status":           common.StatusSuccess,
			"response_code":    strconv.Itoa(responseCode),
		}).Observe(latency.Seconds())
		return
	}
}
