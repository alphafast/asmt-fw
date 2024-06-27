package fcm

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
	"github.com/rs/zerolog"
)

type NotiFCMAdapter struct {
}

func New() *NotiFCMAdapter {
	return &NotiFCMAdapter{}
}

func (n *NotiFCMAdapter) Send(ctx context.Context, req model.NotiRequest) error {
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("[NotiFCMAdapter.Send]: send push notification via FCM")

	return nil
}
