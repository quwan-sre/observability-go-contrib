package sarama

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/quwan-sre/observability-go-contrib/metrics/common"
)

type MetricsInterceptor struct {
	brokers string
}

func NewInterceptor(brokers []string) *MetricsInterceptor {
	brokersTmp := make([]string, len(brokers), len(brokers))
	for i := range brokersTmp {
		brokersTmp[i] = brokers[i]
	}

	sort.Slice(brokersTmp, func(i, j int) bool {
		if brokersTmp[i] < brokersTmp[j] {
			return true
		}
		return false
	})

	return &MetricsInterceptor{
		brokers: strings.Join(brokersTmp, ";"),
	}
}

func (m *MetricsInterceptor) OnSend(msg *sarama.ProducerMessage) {
	common.DefaultMQSendMsgMetric.With(prometheus.Labels{
		"sdk":          common.MQSDKSarama,
		"mq_type":      common.MQTypeKafka,
		"mq_host":      m.brokers,
		"mq_topic":     msg.Topic,
		"mq_partition": common.MQPartitionUnknown,
	}).Inc()
}

func (m *MetricsInterceptor) OnConsume(msg *sarama.ConsumerMessage) {
	common.DefaultMQReceiveMsgMetric.With(prometheus.Labels{
		"sdk":          common.MQSDKSarama,
		"mq_type":      common.MQTypeKafka,
		"mq_host":      m.brokers,
		"mq_topic":     msg.Topic,
		"mq_partition": strconv.FormatInt(int64(msg.Partition), 10),
	}).Inc()
}
