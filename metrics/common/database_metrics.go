package common

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultDatabaseSendRequestMetricName = "apm_database_send_request_duration_milliseconds"
)

var (
	DefaultDatabaseSendRequestMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:                            DefaultDatabaseSendRequestMetricName,
		Buckets:                         []float64{1, 2.5, 5, 10, 20, 50, 100, 500, 1000, 2500, 5000, 7500, 10000},
		NativeHistogramBucketFactor:     1.4,
		NativeHistogramZeroThreshold:    1,
		NativeHistogramMinResetDuration: 5 * time.Minute,
		NativeHistogramMaxZeroThreshold: 50,
		NativeHistogramMaxBucketNumber:  20,
	}, []string{"sdk", "database_type", "database_addr", "response_status", "query_type"})
)
