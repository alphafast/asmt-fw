package noti

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
)

//go:generate mockgen -destination=./mock/repository.go -package=noti_mock github.com/alphafast/asmt-fw/libs/domain/noti NotiRepository

type NotiRepository interface {
	GetNotifyResultsByReqID(ctx context.Context, reqID string) ([]model.NotiResult, error)
	FindUserNotification(ctx context.Context, userID string) (*model.NotiUser, error)
	UpsertNotifyResult(ctx context.Context, result model.NotiResult) error
}
