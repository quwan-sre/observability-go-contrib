package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"time"
)

func init() {
	prometheus.MustRegister(
		// rpc metrics
		DefaultRPCReceiveRequestMetric,
		DefaultRPCSendRequestMetric,

		// cache metrics
		DefaultCacheRequestMetric,
	)

	go func() {
		for {
			// reset all metrics every 25-35 minutes
			time.Sleep(time.Duration(25+rand.Intn(10)) * time.Minute)

			// rpc metrics
			DefaultRPCReceiveRequestMetric.Reset()
			DefaultRPCSendRequestMetric.Reset()

			// cache metrics
			DefaultCacheRequestMetric.Reset()
		}
	}()
}
