package gin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	grpc "google.golang.org/grpc/codes"

	"github.com/quwan-sre/observability-go-contrib/metrics/common"
)

func NewMetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// before
		startTime := time.Now()
		defer func() {
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
				"sdk":              common.RPCSDKGin,
				"request_protocol": common.RPCProtocolHTTP,
				"endpoint":         endpoint,
				"rpc_status_code":  strconv.Itoa(int(grpc.OK)),
				"http_status_code": strconv.Itoa(responseCode),
			}).Observe(latency.Seconds())
		}()
		// execute
		ctx.Next()

		// after

		return
	}
}
