package utils

import (
	"context"

	"github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"

	kafkaUtil "github.com/alphafast/asmt-fw/libs/utils/kafka"
)

type Consumer struct {
	ctx      context.Context
	host     string
	groupId  string
	topic    string
	maxBytes int
}

type ConsumerConfig struct {
	Context context.Context
	Host    string
	GroupId string
	Topic   string
}

type ConsumerHandler func(m kafkaUtil.KafkaMessage) error

func NewConsumer(c ConsumerConfig) *Consumer {
	return &Consumer{
		ctx:      c.Context,
		host:     c.Host,
		groupId:  c.GroupId,
		topic:    c.Topic,
		maxBytes: 10e6,
	}
}

func (c *Consumer) Consume(handler ConsumerHandler) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{c.host},
		GroupID:  c.groupId,
		Topic:    c.topic,
		MaxBytes: c.maxBytes,
	})

	for {
		m, err := r.FetchMessage(c.ctx)
		if err != nil {
			log.Panic().Err(err).Msg("[Consumer.Consume]: error occur while fetch message")
		}

		if err := handler(m); err != nil {
			log.Panic().Err(err).Msg("[Consumer.Consume]: error occur while handling message")
		}

		if err := r.CommitMessages(c.ctx, m); err != nil {
			log.Panic().Err(err).Msg("[Consumer.Consume]: error occur while commit message")
		}
	}
}
