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
			endpoint := common.RPCUnknownString
			host := common.RPCUnknownString
			responseCode := -1

			if ctx != nil && ctx.Request != nil {
				if ctx.Request.Host != "" {
					host = ctx.Request.Host
				}
				if ctx.Request.URL != nil && ctx.Request.URL.Path != "" {
					endpoint = ctx.Request.URL.Path
				}
			}

			if ctx != nil && ctx.Writer != nil {
				responseCode = ctx.Writer.Status()
			}

			common.DefaultRPCReceiveRequestMetric.With(prometheus.Labels{
				"sdk":                  common.RPCSDKGin,
				"request_protocol":     common.RPCProtocolHTTP,
				"request_target":       host,
				"request_path":         endpoint,
				"grpc_response_status": strconv.Itoa(int(grpc.OK)),
				"response_code":        strconv.Itoa(responseCode),
			}).Observe(latency.Seconds() * 1000)
		}()
		// execute
		ctx.Next()

		// after

		return
	}
}
