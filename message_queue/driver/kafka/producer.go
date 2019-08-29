package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/kilingzhang/go-utils/message_queue"
	"github.com/kilingzhang/go-utils/message_queue/driver"
)

var (
	product sarama.SyncProducer
)

type KakfaDriver struct{}

type MySyncProducer struct {
}

func init() {
	messageQueue.Register("kafka", &KakfaDriver{})
}

func (d *KakfaDriver) NewSyncProducer(addrs []string, config interface{}) (driver.SyncProducer, error) {

	if value, ok := config.(*sarama.Config); ok {

		var err error
		product, err = sarama.NewSyncProducer(addrs, value)

		if err != nil {
			return nil, err
		}

		return &MySyncProducer{}, nil
	} else {
		panic("config.(*sarama.Config) fail")
	}
}

func (sp *MySyncProducer) SendMessage(message interface{}, priority int32) (partition int32, offset int64, err error) {
	return product.SendMessage(message.(*sarama.ProducerMessage))
}

func (sp *MySyncProducer) SendMessages(messages []interface{}, priority int32) (err error) {

	var messageArr []*sarama.ProducerMessage

	for _, message := range messages {
		messageArr = append(messageArr, message.(*sarama.ProducerMessage))
	}

	return product.SendMessages(messageArr)
}

func (sp *MySyncProducer) Close() (err error) {
	return product.Close()
}
