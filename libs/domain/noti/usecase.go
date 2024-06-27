package noti

import (
	"context"

	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
)

type RequestNotificationBySourceEvent struct {
	SourceEvent   model.SourceEvent                `json:"sourceEvent"`
	ItemShipped   *ItemShippedInput                `json:"itemShipped,omitempty"`
	ChatMessage   *ChatInput                       `json:"chatMessage,omitempty"`
	BuyerPurchase *BuyerPurchaseInput              `json:"buyerPurchase,omitempty"`
	RemindToPay   *RemindPurchasePendingOrderInput `json:"remindToPay,omitempty"`
}

type ItemShippedInput struct {
	// any item shipped relate input here
	BuyerUserID string   `json:"buyerUserId"`
	Items       []string `json:"items"`
}

type ChatInput struct {
	// any chat relate input here
	SellerUserID string `json:"sellerUserId"`
	Messages     string `json:"messages"`
}

type BuyerPurchaseInput struct {
	// any buyer purchase relate input here
	SellerUserID string `json:"sellerUserId"`
	OrderID      string `json:"orderId"`
}

type RemindPurchasePendingOrderInput struct {
	// any remind purchase pending order relate input here
	BuyerUserID string `json:"buyerUserId"`
	OrderID     string `json:"orderId"`
}

type NotiUseCase interface {
	BuildNotiRequestBySourceEvent(ctx context.Context, req RequestNotificationBySourceEvent) ([]model.NotiRequest, error)
	Notify(ctx context.Context, reqs []model.NotiRequest) ([]model.NotiResult, error)
	SaveNotifyResult(ctx context.Context, res model.NotiResult) error
	GetNotifyResultsByReqID(ctx context.Context, reqID string) ([]model.NotiResult, error)
}
