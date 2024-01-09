package common

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	DefaultMQReceiveMsgMetricName = "apm_mq_receive_msg_count"
	DefaultMQSendMsgMetricName    = "apm_mq_send_msg_count"
)

var (
	DefaultMQReceiveMsgMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: DefaultCacheRequestMetricName,
	}, []string{"sdk", "mq_type", "mq_host", "mq_topic", "mq_partition"})
	DefaultMQSendMsgMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: DefaultMQSendMsgMetricName,
	}, []string{"sdk", "mq_type", "mq_host", "mq_topic", "mq_partition"})
)
