package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultCacheRequestMetricName = "apm_cache_send_request_duration_milliseconds"
)

var (
	DefaultCacheRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultCacheRequestMetricName,
		Buckets:                         []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 3000, 5000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.1,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 10,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "cache_type", "cache_addr", "command", "response_status"})
)
