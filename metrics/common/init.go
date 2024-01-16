package common

import (
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

	// init LRU cache for metrics
	{
		NewRPCSendRequestCache()
		NewRPCReceiveRequestCache()
	}
}
