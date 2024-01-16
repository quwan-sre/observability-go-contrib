package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultRPCReceiveRequestMetricName = "apm_rpc_receive_request_duration_milliseconds"
	DefaultRPCSendRequestMetricName    = "apm_rpc_send_request_duration_milliseconds"
)

var (
	DefaultRPCReceiveRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCReceiveRequestMetricName,
		Buckets:                         []float64{0.5, 1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 30000, 60000, 300000, 600000, 1800000, 3600000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.5,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 50,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "request_target", "request_path", "grpc_response_status", "response_code"})

	DefaultRPCSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultRPCSendRequestMetricName,
		Buckets:                         []float64{0.5, 1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 30000, 60000, 300000, 600000, 1800000, 3600000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.5,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 50,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "request_protocol", "request_target", "request_path", "grpc_response_status", "response_code"})
)
