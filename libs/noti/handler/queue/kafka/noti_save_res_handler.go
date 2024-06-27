package kafka

import (
	"context"
	"encoding/json"

	"github.com/alphafast/asmt-fw/libs/utils/kafka"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/alphafast/asmt-fw/libs/domain/noti"
	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
)

type SaveNotiResultConsumer struct {
	consumerCtx context.Context
	notiUseCase noti.NotiUseCase
	notiRepo    noti.NotiRepository
}

func NewSaveNotiResultConsumerHandler(ctx context.Context, useCase noti.NotiUseCase, repo noti.NotiRepository) *SaveNotiResultConsumer {
	return &SaveNotiResultConsumer{
		consumerCtx: ctx,
		notiUseCase: useCase,
		notiRepo:    repo,
	}
}

func (c *SaveNotiResultConsumer) HandleNotiSaveResult(m kafka.KafkaMessage) error {
	logger := zerolog.Ctx(c.consumerCtx)

	var reqJson model.NotiResult
	if err := json.Unmarshal(m.Value, &reqJson); err != nil {
		logger.Error().Err(err).Msg("[SaveNotiResultConsumer.HandleNotiSaveResult]: failed to unmarshal message")

		return err
	}

	err := c.notiUseCase.SaveNotifyResult(c.consumerCtx, reqJson)
	if err != nil {
		return errors.Wrap(err, "[SaveNotiResultConsumer.HandleNotiSaveResult]: failed to save notify result")
	}

	logger.Info().Msg("[SaveNotiResultConsumer.HandleNotiSaveResult]: notification result saved successfully")

	return nil
}
