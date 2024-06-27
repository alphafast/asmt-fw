package sengrid

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	"github.com/rs/zerolog"
)

type SendGridNotiAdapter struct {
}

func New() *SendGridNotiAdapter {
	return &SendGridNotiAdapter{}
}

func (s *SendGridNotiAdapter) Send(ctx context.Context, req model.NotiRequest) error {
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("[SendGridNotiAdapter.Send]: sending email notification via SendGrid")

	return nil
}
