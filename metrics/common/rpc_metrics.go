package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultRPCReceiveRequestMetricName = "apm_rpc_receive_request_duration_seconds"
	DefaultRPCSendRequestMetricName    = "apm_rpc_send_request_duration_seconds"
)

var (
	DefaultRPCReceiveRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCReceiveRequestMetricName,
		Buckets:                         []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.001,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 0.05,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "endpoint", "rpc_status_code", "http_status_code"})

	DefaultRPCSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCSendRequestMetricName,
		Buckets:                         []float64{0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.001,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 0.05,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "endpoint", "rpc_status_code", "http_status_code"})
)
