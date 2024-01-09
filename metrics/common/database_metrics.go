package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	DefaultDatabaseSendRequestMetricName = "apm_database_send_request_duration_seconds"
)

var (
	DefaultDatabaseSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultDatabaseSendRequestMetricName,
		Buckets:                         []float64{0.001, 0.0025, 0.005, 0.01, 0.02, 0.05, 0.1, 0.5, 1, 2.5, 5, 7.5, 10},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    0.001,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 0.05,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "database_type", "database_addr", "response_status", "query_type"})
)
