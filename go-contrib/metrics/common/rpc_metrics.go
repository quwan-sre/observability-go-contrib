package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	DefaultRPCReceiveRequestMetricName = "apm_rpc_receive_request_duration_seconds"
	DefaultRPCSendRequestMetricName    = "apm_rpc_send_request_duration_seconds"
)

var (
	DefaultRPCReceiveRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCReceiveRequestMetricName,
		Buckets:                         prometheus.DefBuckets,
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.01,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 0.05,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "endpoint", "status", "response_code"})

	DefaultRPCSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCSendRequestMetricName,
		Buckets:                         prometheus.DefBuckets,
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.01,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 0.05,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "endpoint", "status", "response_code"})
)

func init() {
	prometheus.MustRegister(
		DefaultRPCReceiveRequestMetric,
		DefaultRPCSendRequestMetric,
	)
}
