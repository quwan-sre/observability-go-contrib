package metrics

import (
	"fmt"
	"github.com/Shopify/sarama"
	metrics "github.com/quwan-sre/observability-go-contrib/metrics/sarama"
	"testing"
	"time"
)

var (
	kafkaClient sarama.Client
)

func initKafkaClient() {
	brokers := []string{"127.0.0.1:9092"}
	var err error
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V0_11_0_0
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Producer.Interceptors = append(cfg.Producer.Interceptors, metrics.NewInterceptor(brokers))
	cfg.Consumer.Interceptors = append(cfg.Consumer.Interceptors, metrics.NewInterceptor(brokers))
	kafkaClient, err = sarama.NewClient(brokers, cfg)
	if err != nil {
		panic(fmt.Sprintf("create kafka client failed: %v", err))
	}
}

func TestKafkaProducerAndConsumer(t *testing.T) {
	initKafkaClient()

	producer, _ := sarama.NewSyncProducerFromClient(kafkaClient)
	consumer, _ := sarama.NewConsumerFromClient(kafkaClient)

	msg := &sarama.ProducerMessage{
		Topic: "e2e_test_topic",
		Value: sarama.StringEncoder(time.Now().String()),
	}

	partition, offset, err := producer.SendMessage(msg)
	fmt.Printf("producer sendmessage, partition: %d, offset: %d, err: %v\n", partition, offset, err)
	consume, err := consumer.ConsumePartition("e2e_test_topic", partition, offset)
	if err != nil {
		t.Fatal(err)
	}
	if message := <-consume.Messages(); message != nil {
		fmt.Println(string(message.Value))
	}
}
