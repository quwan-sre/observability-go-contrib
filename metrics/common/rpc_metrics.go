package common

import (
	"time"

	"github.com/bluele/gcache"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultRPCReceiveRequestMetricName = "apm_rpc_receive_request_duration_milliseconds"
	DefaultRPCSendRequestMetricName    = "apm_rpc_send_request_duration_milliseconds"
)

const (
	// MaxRequestPathCount represent how many request_path could be stored in metrics vec.
	// It's set to avoid high cardinality.
	MaxRequestPathCount = 1000
	// MaxIdleTime represent the TTL for a idle request_path. When a request_path is not
	// touched for MaxIdleTime, it (related Vec) will be deleted.
	MaxIdleTime = 6 * time.Hour
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

var (
	LRUCacheRPCReceiveRequestMetric gcache.Cache
	LRUCacheRPCSendRequestMetric    gcache.Cache
)

func NewRPCSendRequestCache() {
	cb := gcache.New(MaxRequestPathCount).LRU()

	cb = cb.EvictedFunc(func(k, v interface{}) {
		if requestPath, ok := k.(string); ok {
			DefaultRPCSendRequestMetric.DeletePartialMatch(prometheus.Labels{
				"request_path": requestPath,
			})
		}
	})

	LRUCacheRPCSendRequestMetric = cb.Build()

	go func() {
		for {
			time.Sleep(MaxIdleTime)
			allRequestPath := LRUCacheRPCSendRequestMetric.GetALL(false)
			validRequestPath := LRUCacheRPCSendRequestMetric.GetALL(true)
			for k, _ := range allRequestPath {
				if _, ok := validRequestPath[k]; !ok {
					LRUCacheRPCSendRequestMetric.Remove(k)
				}
			}
		}
	}()
}

func NewRPCReceiveRequestCache() {
	cb := gcache.New(MaxRequestPathCount).LRU()

	cb = cb.EvictedFunc(func(k, v interface{}) {
		if requestPath, ok := k.(string); ok {
			DefaultRPCReceiveRequestMetric.DeletePartialMatch(prometheus.Labels{
				"request_path": requestPath,
			})
		}
	})

	LRUCacheRPCReceiveRequestMetric = cb.Build()

	go func() {
		for {
			time.Sleep(MaxIdleTime)
			allRequestPath := LRUCacheRPCReceiveRequestMetric.GetALL(false)
			validRequestPath := LRUCacheRPCReceiveRequestMetric.GetALL(true)
			for k, _ := range allRequestPath {
				if _, ok := validRequestPath[k]; !ok {
					LRUCacheRPCReceiveRequestMetric.Remove(k)
				}
			}
		}
	}()
}
