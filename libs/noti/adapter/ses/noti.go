package ses

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	"github.com/rs/zerolog"
)

type SesNotiAdapter struct {
}

func New() *SesNotiAdapter {
	return &SesNotiAdapter{}
}

func (s *SesNotiAdapter) Send(ctx context.Context, req model.NotiRequest) error {
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("[SesNotiAdapter.Send]: sending email notification via SES")

	return nil
}
