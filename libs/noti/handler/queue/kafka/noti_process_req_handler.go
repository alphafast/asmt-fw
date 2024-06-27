package kafka

import (
	"context"
	"encoding/json"

	"github.com/alphafast/asmt-fw/libs/domain/noti"
	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	"github.com/alphafast/asmt-fw/libs/utils/kafka"
	"github.com/alphafast/asmt-fw/libs/utils/kafka/producer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type ProcessNotiRequestHandlerOps struct {
	ProcessedTopic string
}

type ProcessNotiRequestHandlerDeps struct {
	NotiUseCase noti.NotiUseCase
	Producer    *producer.Producer
}

type ProcessNotiRequestHandler struct {
	consumerCtx context.Context
	notiUseCase noti.NotiUseCase
	producer    *producer.Producer

	processedTopic string
}

func NewProcessNotiRequestHandler(ctx context.Context, d ProcessNotiRequestHandlerDeps, c ProcessNotiRequestHandlerOps) *ProcessNotiRequestHandler {
	return &ProcessNotiRequestHandler{
		consumerCtx: ctx,
		notiUseCase: d.NotiUseCase,
		producer:    d.Producer,

		processedTopic: c.ProcessedTopic,
	}
}

func (c *ProcessNotiRequestHandler) HandleNotiRequest(m kafka.KafkaMessage) error {
	logger := zerolog.Ctx(c.consumerCtx)

	var reqJson model.NotiRequest
	if err := json.Unmarshal(m.Value, &reqJson); err != nil {
		logger.Error().Err(err).Msg("[ProcessNotiRequestHandler.HandleNotiRequest]: failed to unmarshal message")

		return err
	}

	ress, err := c.notiUseCase.Notify(c.consumerCtx, []model.NotiRequest{reqJson})
	if err != nil {
		return errors.Wrap(err, "[ProcessNotiRequestHandler.HandleNotiRequest]: failed to process notification request")
	}

	tobeProduceMsgs := []producer.Message{}
	for _, res := range ress {
		tobeProduceMsgs = append(tobeProduceMsgs, producer.Message{
			Topic:   c.processedTopic,
			Key:     res.ID,
			Message: res,
		})
	}

	if err := c.producer.Produce(c.consumerCtx, tobeProduceMsgs...); err != nil {
		return errors.Wrap(err, "[ProcessNotiRequestHandler.HandleNotiRequest]: failed to produce messages")
	}

	logger.Info().Msgf("[ProcessNotiRequestHandler.HandleNotiRequest]: %d notification(s) processed successfully", len(ress))

	return nil
}
