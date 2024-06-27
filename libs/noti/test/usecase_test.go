package test_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bsm/ginkgo/v2"
	"github.com/bsm/gomega"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"

	noti "github.com/alphafast/asmt-fw/libs/domain/noti"
	notiMock "github.com/alphafast/asmt-fw/libs/domain/noti/mock"
	notiModel "github.com/alphafast/asmt-fw/libs/domain/noti/model"

	uc "github.com/alphafast/asmt-fw/libs/noti/usecase"
)

func TestNotificationUseCase(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Notification Test Suite")
}

var _ = ginkgo.Describe("NotiUseCase", func() {
	var (
		mockCtrl *gomock.Controller

		mockNotiAdapter *notiMock.MockNotiAdapter
		mockNotiRepo    *notiMock.MockNotiRepository

		usecase noti.NotiUseCase
		ctx     context.Context
	)

	ginkgo.BeforeEach(func() {
		mockCtrl = gomock.NewController(ginkgo.GinkgoT())
		mockNotiAdapter = notiMock.NewMockNotiAdapter(mockCtrl)
		mockNotiRepo = notiMock.NewMockNotiRepository(mockCtrl)
		ctx = context.Background()

		deps := uc.NotiUsecaseDeps{
			NotiRepo:    mockNotiRepo,
			NotiAdapter: mockNotiAdapter,
		}
		conf := uc.NotiUsecaseConf{
			NotiChannelByEventSource: notiModel.DefaultNotiChannelBySourceEvent,
		}
		usecase = uc.New(deps, conf)
	})

	ginkgo.Describe("Notify", func() {
		ginkgo.Context("When Notify is called", func() {
			ginkgo.It("should return results", func() {
				reqs := []notiModel.NotiRequest{
					{
						ID:           "some-id-1",
						ReqID:        "some-req-id-1",
						SourceEvent:  notiModel.ItemShippedNotification,
						NotiType:     "email",
						EmailPayload: &notiModel.EmailPayload{},
					},
					{
						ID:          "some-id-2",
						ReqID:       "some-req-id-2",
						SourceEvent: notiModel.ItemShippedNotification,
						NotiType:    "push",
						PushPayload: &notiModel.PushPayload{},
					},
				}
				expectedResults := []notiModel.NotiResult{
					{
						ID:        "some-id-1",
						ReqID:     "some-req-id-1",
						IsSuccess: true,
						Reason:    "",
					},
					{
						ID:        "some-id-2",
						ReqID:     "some-req-id-2",
						IsSuccess: true,
						Reason:    "",
					},
				}

				mockNotiAdapter.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil).Times(2)

				results, err := usecase.Notify(ctx, reqs)

				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(results).To(gomega.HaveLen(2))
				gomega.Expect(results).To(gomega.ConsistOf(expectedResults))
			})
		})
	})

	ginkgo.Describe("SaveNotifyResult", func() {
		ginkgo.Context("When SaveNotifyResult is called", func() {
			ginkgo.It("should return nil", func() {
				res := notiModel.NotiResult{
					ID:        "some-id",
					ReqID:     "some-req-id",
					IsSuccess: true,
					Reason:    "",
				}

				mockNotiRepo.EXPECT().UpsertNotifyResult(gomock.Any(), res).Return(nil)

				err := usecase.SaveNotifyResult(ctx, res)

				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When SaveNotifyResult is called with error", func() {
			ginkgo.It("should return error", func() {
				res := notiModel.NotiResult{
					ID:        "some-id",
					ReqID:     "some-req-id",
					IsSuccess: true,
					Reason:    "",
				}

				mockNotiRepo.EXPECT().UpsertNotifyResult(gomock.Any(), res).Return(errors.New("some-error"))

				err := usecase.SaveNotifyResult(ctx, res)

				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("RequestNotification", func() {
		ginkgo.Context("When RequestNotification is called with any kind of event and user not have push channel", func() {
			ginkgo.It("should return error", func() {
				reqByEventSource := noti.RequestNotificationBySourceEvent{
					SourceEvent: notiModel.ItemShippedNotification,
					ItemShipped: &noti.ItemShippedInput{
						BuyerUserID: "some-buyer-id",
						Items:       []string{"item-1", "item-2"},
					},
				}

				mockNotiRepo.EXPECT().FindUserNotification(gomock.Any(), reqByEventSource.ItemShipped.BuyerUserID).Return(nil, errors.New("some error")).Times(1)

				reqs, err := usecase.BuildNotiRequestBySourceEvent(ctx, reqByEventSource)

				gomega.Expect(err).NotTo(gomega.BeNil())
				gomega.Expect(reqs).To(gomega.BeNil())

			})
		})

		ginkgo.Context("When RequestNotification is called with ItemShippedNotification", func() {
			ginkgo.It("should return noti requests", func() {
				mockDeviceToken := "some-device-token"

				reqByEventSource := noti.RequestNotificationBySourceEvent{
					SourceEvent: notiModel.ItemShippedNotification,
					ItemShipped: &noti.ItemShippedInput{
						BuyerUserID: "some-buyer-id",
						Items:       []string{"item-1", "item-2"},
					},
				}

				expectedReqs := []notiModel.NotiRequest{
					{
						SourceEvent: notiModel.ItemShippedNotification,
						NotiType:    notiModel.PushType,
						PushPayload: &notiModel.PushPayload{
							Title: "Item Shipped",
							Body:  "Your item has been shipped",
							Token: mockDeviceToken,
						},
					},
				}

				mockNotiRepo.EXPECT().FindUserNotification(gomock.Any(), reqByEventSource.ItemShipped.BuyerUserID).Return(&notiModel.NotiUser{
					Channels: []notiModel.NotiUserNotiChannel{
						{
							NotiType: notiModel.PushType,
							PushChannelPayload: &notiModel.NotiUserPushChannelPayload{
								Token: mockDeviceToken,
							},
						},
					},
				}, nil).Times(1)

				reqs, err := usecase.BuildNotiRequestBySourceEvent(ctx, reqByEventSource)

				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(reqs).To(gomega.HaveLen(1))
				gomega.Expect(reqs[0].SourceEvent).To(gomega.Equal((expectedReqs[0].SourceEvent)))
				gomega.Expect(reqs[0].NotiType).To(gomega.Equal((expectedReqs[0].NotiType)))
				gomega.Expect(reqs[0].PushPayload).To(gomega.Equal((expectedReqs[0].PushPayload)))
			})
		})

		ginkgo.Context("When RequestNotification is called with ChatMessageNotification", func() {
			ginkgo.It("should return noti requests", func() {
				mockDeviceToken := "some-device-token"
				mockEmail := "some-email"

				reqByEventSource := noti.RequestNotificationBySourceEvent{
					SourceEvent: notiModel.ChatMessageNotification,
					ChatMessage: &noti.ChatInput{
						SellerUserID: "some-seller-id",
						Messages:     "some-messages",
					},
				}

				mockNotiRepo.EXPECT().FindUserNotification(gomock.Any(), reqByEventSource.ChatMessage.SellerUserID).Return(&notiModel.NotiUser{
					Channels: []notiModel.NotiUserNotiChannel{
						{
							NotiType: notiModel.PushType,
							PushChannelPayload: &notiModel.NotiUserPushChannelPayload{
								Token: mockDeviceToken,
							},
						},
						{
							NotiType: notiModel.EmailType,
							EmailChannelPayload: &notiModel.NotiUserEmailChannelPayload{
								EmailAddress: mockEmail,
							},
						},
					},
				}, nil).Times(1)

				ress, err := usecase.BuildNotiRequestBySourceEvent(ctx, reqByEventSource)

				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(ress).To(gomega.HaveLen(2))
				forEmail, ok := lo.Find(ress, func(res notiModel.NotiRequest) bool {
					return res.NotiType == notiModel.EmailType
				})
				gomega.Expect(ok).To(gomega.BeTrue())
				gomega.Expect(forEmail.SourceEvent).To(gomega.Equal(notiModel.ChatMessageNotification))
				gomega.Expect(forEmail.NotiType).To(gomega.Equal(notiModel.EmailType))
				gomega.Expect(forEmail.EmailPayload).To(gomega.Equal(&notiModel.EmailPayload{
					Subject: "New Message",
					Body:    reqByEventSource.ChatMessage.Messages,
					To:      mockEmail,
				}))

				forPush, ok := lo.Find(ress, func(res notiModel.NotiRequest) bool {
					return res.NotiType == notiModel.PushType
				})
				gomega.Expect(ok).To(gomega.BeTrue())
				gomega.Expect(forPush.SourceEvent).To(gomega.Equal(notiModel.ChatMessageNotification))
				gomega.Expect(forPush.NotiType).To(gomega.Equal(notiModel.PushType))
				gomega.Expect(forPush.PushPayload).To(gomega.Equal(&notiModel.PushPayload{
					Title: "New Message",
					Body:  reqByEventSource.ChatMessage.Messages,
					Token: mockDeviceToken,
				}))
				gomega.Expect(ress[0].ReqID).To(gomega.Equal(ress[1].ReqID))
				gomega.Expect(ress[0].ID).NotTo(gomega.Equal(ress[1].ID))
			})
		})

		ginkgo.Context("When RequestNotification is called with BuyerPurchaseNotification", func() {
			ginkgo.It("should return noti requests", func() {
				mockDeviceToken := "some-device-token"
				mockEmail := "some-email"

				reqByEventSource := noti.RequestNotificationBySourceEvent{
					SourceEvent: notiModel.BuyerPurchaseNotification,
					BuyerPurchase: &noti.BuyerPurchaseInput{
						SellerUserID: "some-seller-id",
						OrderID:      "some-order-id",
					},
				}

				mockNotiRepo.EXPECT().FindUserNotification(gomock.Any(), reqByEventSource.BuyerPurchase.SellerUserID).Return(&notiModel.NotiUser{
					Channels: []notiModel.NotiUserNotiChannel{
						{
							NotiType: notiModel.PushType,
							PushChannelPayload: &notiModel.NotiUserPushChannelPayload{
								Token: mockDeviceToken,
							},
						},
						{
							NotiType: notiModel.EmailType,
							EmailChannelPayload: &notiModel.NotiUserEmailChannelPayload{
								EmailAddress: mockEmail,
							},
						},
					},
				}, nil).Times(1)

				ress, err := usecase.BuildNotiRequestBySourceEvent(ctx, reqByEventSource)

				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(ress).To(gomega.HaveLen(2))
				forEmail, ok := lo.Find(ress, func(res notiModel.NotiRequest) bool {
					return res.NotiType == notiModel.EmailType
				})
				gomega.Expect(ok).To(gomega.BeTrue())
				gomega.Expect(forEmail.SourceEvent).To(gomega.Equal(notiModel.BuyerPurchaseNotification))
				gomega.Expect(forEmail.NotiType).To(gomega.Equal(notiModel.EmailType))
				gomega.Expect(forEmail.EmailPayload).To(gomega.Equal(&notiModel.EmailPayload{
					Subject: "Buyer Purchase updated",
					Body:    "Buyer have purchased order some-order-id, more information <a href=\"someweb.com/order/some-order-id\">click</a>",
					To:      mockEmail,
				}))
				forPush, ok := lo.Find(ress, func(res notiModel.NotiRequest) bool {
					return res.NotiType == notiModel.PushType
				})
				gomega.Expect(ok).To(gomega.BeTrue())
				gomega.Expect(forPush.SourceEvent).To(gomega.Equal(notiModel.BuyerPurchaseNotification))
				gomega.Expect(forPush.NotiType).To(gomega.Equal(notiModel.PushType))
				gomega.Expect(forPush.PushPayload).To(gomega.Equal(&notiModel.PushPayload{
					Title: "Buyer Purchase updated",
					Body:  "Buyer have purchased order some-order-id",
					Token: mockDeviceToken,
				}))
				gomega.Expect(ress[0].ReqID).To(gomega.Equal(ress[1].ReqID))
				gomega.Expect(ress[0].ID).NotTo(gomega.Equal(ress[1].ID))
			})
		})

		ginkgo.Context("When RequestNotification is called with RemindPurchasePendingOrderNotification", func() {
			ginkgo.It("should return noti requests", func() {
				mockDeviceToken := "some device token"

				reqByEventSource := noti.RequestNotificationBySourceEvent{
					SourceEvent: notiModel.RemindToPayOrderNotification,
					RemindToPay: &noti.RemindPurchasePendingOrderInput{
						BuyerUserID: "some-buyer-id",
						OrderID:     "some-order-id",
					},
				}

				mockNotiRepo.EXPECT().FindUserNotification(gomock.Any(), reqByEventSource.RemindToPay.BuyerUserID).Return(&notiModel.NotiUser{
					Channels: []notiModel.NotiUserNotiChannel{
						{
							NotiType: notiModel.PushType,
							PushChannelPayload: &notiModel.NotiUserPushChannelPayload{
								Token: mockDeviceToken,
							},
						},
					},
				}, nil).Times(1)

				ress, err := usecase.BuildNotiRequestBySourceEvent(ctx, reqByEventSource)

				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(ress).To(gomega.HaveLen(1))
				gomega.Expect(ress[0].SourceEvent).To(gomega.Equal(notiModel.RemindToPayOrderNotification))
				gomega.Expect(ress[0].NotiType).To(gomega.Equal(notiModel.PushType))
				gomega.Expect(ress[0].PushPayload).To(gomega.Equal(&notiModel.PushPayload{
					Title: "Remind to pay order",
					Body:  "You have pending order some-order-id, please pay",
					Token: mockDeviceToken,
				}))
			})
		})
	})

	ginkgo.AfterEach(func() {
		mockCtrl.Finish()
	})
})
