package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/kilingzhang/go-utils/message_queue"
	"github.com/kilingzhang/go-utils/message_queue/driver"
)

var (
	syncProduct   sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
)

type KakfaDriver struct{}

type KafkaSyncProducer struct{}

type KafkaConsumerGroup struct{}

func init() {
	messageQueue.Register("kafka", &KakfaDriver{})
}

func (d *KakfaDriver) NewSyncProducer(addrs []string, config interface{}) (driver.SyncProducer, error) {

	if value, ok := config.(*sarama.Config); ok {

		var err error
		syncProduct, err = sarama.NewSyncProducer(addrs, value)

		if err != nil {
			return nil, err
		}

		return &KafkaSyncProducer{}, nil
	} else {
		panic("config.(*sarama.Config) fail")
	}
}

func (d *KakfaDriver) NewConsumerGroup(addrs []string, groupID string, config interface{}) (driver.ConsumerGroup, error) {

	if value, ok := config.(*sarama.Config); ok {

		var err error
		consumerGroup, err = sarama.NewConsumerGroup(addrs, groupID, value)

		if err != nil {
			return nil, err
		}

		return &KafkaConsumerGroup{}, nil
	} else {
		panic("config.(*sarama.Config) fail")
	}
}

func (sp *KafkaSyncProducer) SendMessage(message interface{}, priority int32) (partition int32, offset int64, err error) {

	return syncProduct.SendMessage(message.(*sarama.ProducerMessage))
}

func (sp *KafkaSyncProducer) SendMessages(messages []interface{}, priority int32) (err error) {

	var messageArr []*sarama.ProducerMessage

	for _, message := range messages {
		messageArr = append(messageArr, message.(*sarama.ProducerMessage))
	}

	return syncProduct.SendMessages(messageArr)
}

func (sp *KafkaSyncProducer) Close() (err error) {

	return syncProduct.Close()
}

func (sp *KafkaConsumerGroup) Consume(ctx context.Context, topics []string, handler interface{}) error {

	return consumerGroup.Consume(ctx, topics, handler.(sarama.ConsumerGroupHandler))
}

func (sp *KafkaConsumerGroup) Errors() <-chan error {

	return consumerGroup.Errors()
}

func (sp *KafkaConsumerGroup) Close() error {

	return consumerGroup.Close()
}
