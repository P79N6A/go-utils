package driver

import (
	"context"
)

type Driver interface {
	NewSyncProducer(addrs []string, config interface{}) (SyncProducer, error)
	NewConsumerGroup(addrs []string, groupID string, config interface{}) (ConsumerGroup, error)
}

type SyncProducer interface {
	SendMessage(message interface{}, priority int32) (partition int32, offset int64, err error)

	SendMessages(messages []interface{}, priority int32) (err error)

	Close() (err error)
}

type ConsumerGroup interface {
	Consume(ctx context.Context, topics []string, handler interface{}) error

	Errors() <-chan error

	Close() error
}
