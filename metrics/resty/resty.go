package resty

import (
	"context"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	grpc "google.golang.org/grpc/codes"

	"github.com/quwan-sre/observability-go-contrib/metrics/common"
)

type metricsCtxKey struct{}

func NewMetricsMiddleware(c *resty.Client) {
	if c == nil {
		return
	}

	c.OnBeforeRequest(NewBeforeRequest())
	c.OnAfterResponse(NewAfterResponse())

	return
}

func NewBeforeRequest() func(c *resty.Client, r *resty.Request) error {
	return func(c *resty.Client, r *resty.Request) error {
		ctx := context.WithValue(r.Context(), metricsCtxKey{}, time.Now())
		r.SetContext(ctx)
		return nil
	}
}

func NewAfterResponse() func(c *resty.Client, r *resty.Response) error {
	return func(c *resty.Client, r *resty.Response) error {
		req := r.Request
		timeInterface := req.Context().Value(metricsCtxKey{})
		t, ok := timeInterface.(time.Time)
		if !ok {
			// should never reach here
			return nil
		}

		latency := time.Now().Sub(t)
		endpoint := ""
		if req.RawRequest != nil && req.RawRequest.URL != nil {
			endpoint = req.RawRequest.URL.Path
		}

		common.DefaultRPCSendRequestMetric.With(prometheus.Labels{
			"sdk":              common.RPCSDKResty,
			"request_protocol": common.RPCProtocolHTTP,
			"endpoint":         endpoint,
			"rpc_status":       strconv.Itoa(int(grpc.OK)),
			"response_code":    strconv.Itoa(r.StatusCode()),
		}).Observe(latency.Seconds())
		return nil
	}
}
