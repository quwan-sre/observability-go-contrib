## sarama

sarama 提供了两组 Interface，用于扩展额外逻辑。
```go
// ProducerInterceptor allows you to intercept (and possibly mutate) the records
// received by the producer before they are published to the Kafka cluster.
// https://cwiki.apache.org/confluence/display/KAFKA/KIP-42%3A+Add+Producer+and+Consumer+Interceptors#KIP42:AddProducerandConsumerInterceptors-Motivation
type ProducerInterceptor interface {

	// OnSend is called when the producer message is intercepted. Please avoid
	// modifying the message until it's safe to do so, as this is _not_ a copy
	// of the message.
	OnSend(*ProducerMessage)
}

// ConsumerInterceptor allows you to intercept (and possibly mutate) the records
// received by the consumer before they are sent to the messages channel.
// https://cwiki.apache.org/confluence/display/KAFKA/KIP-42%3A+Add+Producer+and+Consumer+Interceptors#KIP42:AddProducerandConsumerInterceptors-Motivation
type ConsumerInterceptor interface {

	// OnConsume is called when the consumed message is intercepted. Please
	// avoid modifying the message until it's safe to do so, as this is _not_ a
	// copy of the message.
	OnConsume(*ConsumerMessage)
}
```

这两组 Interface 仅允许用户在消息发出前或消费后执行逻辑，因此无法统计耗时。

在 sarama 中可以通过以下代码使用 go-contrib 埋点：
```go
func initKafkaClient() {
	brokers := []string{"127.0.0.1:9092"}
	var err error
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V0_11_0_0
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	
	// 追加自定义的 Interceptor
	cfg.Producer.Interceptors = append(cfg.Producer.Interceptors, metrics.NewInterceptor(brokers))
	cfg.Consumer.Interceptors = append(cfg.Consumer.Interceptors, metrics.NewInterceptor(brokers))
	
	kafkaClient, err = sarama.NewClient(brokers, cfg)
	if err != nil {
		panic(fmt.Sprintf("create kafka client failed: %v", err))
	}
}
```

由于 Interface 中可以获取的信息不足，在初始化时需要提供 Kafka 的相关信息（broker 地址），该信息会作为 Label 填充至相关 Metrics。
