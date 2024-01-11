package common

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(
		// rpc metrics
		DefaultRPCReceiveRequestMetric,
		DefaultRPCSendRequestMetric,

		// cache metrics
		DefaultCacheRequestMetric,

		// mq metrics
		DefaultMQSendMsgMetric,
		DefaultMQReceiveMsgMetric,

		// database metrics
		DefaultDatabaseSendRequestMetric,
	)

	go func() {
		for {
			// reset all metrics every 23-24 hours
			time.Sleep(time.Duration((60*23)+rand.Intn(60)) * time.Minute)

			// rpc metrics
			DefaultRPCReceiveRequestMetric.Reset()
			DefaultRPCSendRequestMetric.Reset()

			// cache metrics
			DefaultCacheRequestMetric.Reset()

			// mq metrics
			DefaultMQReceiveMsgMetric.Reset()
			DefaultMQSendMsgMetric.Reset()

			// database metrics
			DefaultDatabaseSendRequestMetric.Reset()
		}
	}()
}
