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
			// reset all metrics every 25-35 minutes
			time.Sleep(time.Duration(25+rand.Intn(10)) * time.Minute)

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
