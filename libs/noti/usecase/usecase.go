package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/alphafast/asmt-fw/libs/domain/noti"
	"github.com/alphafast/asmt-fw/libs/domain/noti/model"
)

var (
	createUUID = uuid.New
)

type NotiUsecase struct {
	notiRepo    noti.NotiRepository
	notiAdapter noti.NotiAdapter
}

func New(notiRepo noti.NotiRepository, notiAdapter noti.NotiAdapter) *NotiUsecase {
	return &NotiUsecase{
		notiRepo:    notiRepo,
		notiAdapter: notiAdapter,
	}
}

func (uc *NotiUsecase) BuildNotiRequestBySourceEvent(ctx context.Context, req noti.RequestNotificationBySourceEvent) ([]model.NotiRequest, error) {
	var reqs []model.NotiRequest
	var err error

	switch req.SourceEvent {
	case model.ItemShippedNotification:
		reqs, err = uc.buildItemShippedNotification(ctx, createUUID().String(), *req.ItemShipped)
		if err != nil {
			return nil, errors.Wrap(err, "[NotiUsecase.RequestNotification] BuildItemShippedNotification error")
		}
	case model.ChatMessageNotification:
		reqs, err = uc.buildChatMessageNotification(ctx, createUUID().String(), *req.ChatMessage)
		if err != nil {
			return nil, errors.Wrap(err, "[NotiUsecase.RequestNotification] BuildChatMessageNotification error")
		}
	case model.BuyerPurchaseNotification:
		reqs, err = uc.buildBuyerPurchaseNotification(ctx, createUUID().String(), *req.BuyerPurchase)
		if err != nil {
			return nil, errors.Wrap(err, "[NotiUsecase.RequestNotification] BuildBuyerPurchaseNotification error")
		}
	case model.RemindToPayOrderNotification:
		reqs, err = uc.buildRemindToPayOrderNotification(ctx, createUUID().String(), *req.RemindToPay)
		if err != nil {
			return nil, errors.Wrap(err, "[NotiUsecase.RequestNotification] BuildRemindToPayOrderNotification error")
		}
	default:
		return nil, errors.New("[NotiUsecase.RequestNotification] Invalid source event")
	}

	return reqs, nil
}

func (uc *NotiUsecase) Notify(ctx context.Context, reqs []model.NotiRequest) ([]model.NotiResult, error) {
	results := []model.NotiResult{}
	g := errgroup.Group{}

	for _, req := range reqs {
		eachReq := req
		g.Go(func() error {
			err := uc.notiAdapter.Send(ctx, eachReq)
			if err != nil {
				results = append(results, model.NotiResult{
					ID:        eachReq.ID,
					ReqID:     eachReq.ReqID,
					IsSuccess: false,
					Reason:    err.Error(),
				})

				return nil
			}

			results = append(results, model.NotiResult{
				ID:        eachReq.ID,
				ReqID:     eachReq.ReqID,
				IsSuccess: true,
				Reason:    "",
			})

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.Wrap(err, "[NotiUsecase.Notify] error occurred")
	}

	return results, nil
}

func (uc *NotiUsecase) SaveNotifyResult(ctx context.Context, res model.NotiResult) error {
	err := uc.notiRepo.UpsertNotifyResult(ctx, res)
	if err != nil {
		return errors.Wrap(err, "[NotiUsecase.SaveNotifyResult] error occurred")
	}

	return nil
}

func (uc *NotiUsecase) GetNotifyResultsByReqID(ctx context.Context, reqID string) ([]model.NotiResult, error) {
	results, err := uc.notiRepo.GetNotifyResultsByReqID(ctx, reqID)
	if err != nil {
		return nil, errors.Wrap(err, "[NotiUsecase.GetNotifyResultsByReqID] error occurred")
	}

	return results, nil

}

func (uc *NotiUsecase) buildItemShippedNotification(ctx context.Context, reqID string, input noti.ItemShippedInput) ([]model.NotiRequest, error) {
	targetUserID := input.BuyerUserID
	user, err := uc.notiRepo.FindUserNotification(ctx, targetUserID)
	if err != nil {
		return nil, errors.Wrap(err, "[NotiUsecase.BuildItemShippedNotification] FindUserNotification error")
	}

	chanMap := user.GetTypeChannelMap()
	if _, ok := chanMap[model.PushType]; !ok {
		return nil, errors.New("[NotiUsecase.BuildItemShippedNotification] User not have push channel")
	}

	pushChannel := chanMap[model.PushType].PushChannelPayload

	return []model.NotiRequest{
		{
			ID:          createUUID().String(),
			ReqID:       reqID,
			SourceEvent: model.ItemShippedNotification,
			NotiType:    model.PushType,
			PushPayload: &model.PushPayload{
				Title: "Item Shipped",
				Body:  "Your item has been shipped",
				Token: pushChannel.Token,
			},
		},
	}, nil
}

func (uc *NotiUsecase) buildChatMessageNotification(ctx context.Context, reqID string, input noti.ChatInput) ([]model.NotiRequest, error) {
	targetUserID := input.SellerUserID
	user, err := uc.notiRepo.FindUserNotification(ctx, targetUserID)
	if err != nil {
		return nil, errors.Wrap(err, "[NotiUsecase.BuildChatMessageNotification] FindUserNotification error")
	}

	chanMap := user.GetTypeChannelMap()
	if _, ok := chanMap[model.PushType]; !ok {
		return nil, errors.New("[NotiUsecase.BuildChatMessageNotification] User not have push channel")
	}
	if _, ok := chanMap[model.EmailType]; !ok {
		return nil, errors.New("[NotiUsecase.BuildChatMessageNotification] User not have email channel")
	}
	pushChannel := chanMap[model.PushType].PushChannelPayload
	emailChannel := chanMap[model.EmailType].EmailChannelPayload

	return []model.NotiRequest{
		{
			ID:          createUUID().String(),
			ReqID:       reqID,
			SourceEvent: model.ChatMessageNotification,
			NotiType:    model.PushType,
			PushPayload: &model.PushPayload{
				Title: "New Message",
				Body:  input.Messages,
				Token: pushChannel.Token,
			},
		},
		{
			ID:          createUUID().String(),
			ReqID:       reqID,
			SourceEvent: model.ChatMessageNotification,
			NotiType:    model.EmailType,
			EmailPayload: &model.EmailPayload{
				Subject: "New Message",
				Body:    input.Messages,
				To:      emailChannel.EmailAddress,
			},
		},
	}, nil
}

func (uc *NotiUsecase) buildBuyerPurchaseNotification(ctx context.Context, reqID string, input noti.BuyerPurchaseInput) ([]model.NotiRequest, error) {
	targetUserID := input.SellerUserID
	user, err := uc.notiRepo.FindUserNotification(ctx, targetUserID)
	if err != nil {
		return nil, errors.Wrap(err, "[NotiUsecase.BuildBuyerPurchaseNotification] FindUserNotification error")
	}

	chanMap := user.GetTypeChannelMap()
	if _, ok := chanMap[model.PushType]; !ok {
		return nil, errors.New("[NotiUsecase.BuildBuyerPurchaseNotification] User not have push channel")
	}
	if _, ok := chanMap[model.EmailType]; !ok {
		return nil, errors.New("[NotiUsecase.BuildBuyerPurchaseNotification] User not have email channel")
	}

	pushChannel := chanMap[model.PushType].PushChannelPayload
	emailChannel := chanMap[model.EmailType].EmailChannelPayload

	return []model.NotiRequest{
		{
			ID:          createUUID().String(),
			ReqID:       reqID,
			SourceEvent: model.BuyerPurchaseNotification,
			NotiType:    model.PushType,
			PushPayload: &model.PushPayload{
				Title: "Buyer Purchase updated",
				Body:  fmt.Sprintf("Buyer have purchased order %s", input.OrderID),
				Token: pushChannel.Token,
			},
		},
		{
			ID:          createUUID().String(),
			ReqID:       reqID,
			SourceEvent: model.BuyerPurchaseNotification,
			NotiType:    model.EmailType,
			EmailPayload: &model.EmailPayload{
				Subject: "Buyer Purchase updated",
				Body:    fmt.Sprintf("Buyer have purchased order %s, more information <a href=\"someweb.com/order/%s\">click</a>", input.OrderID, input.OrderID),
				To:      emailChannel.EmailAddress,
			},
		},
	}, nil
}

func (uc *NotiUsecase) buildRemindToPayOrderNotification(ctx context.Context, reqID string, input noti.RemindPurchasePendingOrderInput) ([]model.NotiRequest, error) {
	targetUserID := input.BuyerUserID
	user, err := uc.notiRepo.FindUserNotification(ctx, targetUserID)
	if err != nil {
		return nil, errors.Wrap(err, "[NotiUsecase.BuildRemindToPayOrderNotification] FindUserNotification error")
	}

	chanMap := user.GetTypeChannelMap()
	if _, ok := chanMap[model.PushType]; !ok {
		return nil, errors.New("[NotiUsecase.BuildRemindToPayOrderNotification] User not have push channel")
	}

	pushChannel := chanMap[model.PushType].PushChannelPayload

	return []model.NotiRequest{
		{
			ID:          createUUID().String(),
			ReqID:       reqID,
			SourceEvent: model.RemindToPayOrderNotification,
			NotiType:    model.PushType,
			PushPayload: &model.PushPayload{
				Title: "Remind to pay order",
				Body:  fmt.Sprintf("You have pending order %s, please pay", input.OrderID),
				Token: pushChannel.Token,
			},
		},
	}, nil
}
