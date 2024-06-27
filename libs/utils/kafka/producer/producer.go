package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	host   string
	writer *kafka.Writer
}

func NewProducer(host string) *Producer {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(host),
		Balancer:               &kafka.Murmur2Balancer{},
		AllowAutoTopicCreation: false,
		Async:                  true,
	}

	return &Producer{
		host:   host,
		writer: &writer,
	}
}

type Message struct {
	Topic   string
	Key     string
	Header  []kafka.Header
	Message interface{}
}

func (p *Producer) Produce(ctx context.Context, msg ...Message) error {
	kfkMsgs := []kafka.Message{}
	for _, m := range msg {
		bb, err := json.Marshal(m.Message)
		if err != nil {
			log.Error().Err(err).Msg("[Producer.Produce] error while marshalling message")

			return errors.Wrap(err, "[Producer.Produce] error while marshalling message")
		}

		key := []byte{}
		if m.Key == "" {
			key = []byte(uuid.New().String())
		} else {
			key = []byte(m.Key)
		}

		kfkMsgs = append(kfkMsgs, kafka.Message{
			Topic:   m.Topic,
			Key:     key,
			Headers: m.Header,
			Value:   bb,
		})
	}

	// Produce message to kafka
	err := p.writer.WriteMessages(ctx, kfkMsgs...)
	if err != nil {
		log.Error().Err(err).Msg("[Producer.Produce] error while writing messages to kafka")

		return errors.Wrap(err, "[Producer.Produce] error while writing messages to kafka")
	}

	return nil
}
