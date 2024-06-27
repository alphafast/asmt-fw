package noti

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
)

//go:generate mockgen -destination=./mock/adapter.go -package=noti_mock github.com/alphafast/asmt-fw/libs/domain/noti NotiAdapter

type NotiAdapter interface {
	Send(ctx context.Context, req model.NotiRequest) error
}
