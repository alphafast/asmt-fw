package model

type SourceEvent string

type NotiType string

const (
	ItemShippedNotification      SourceEvent = "item-shipped"
	ChatMessageNotification      SourceEvent = "chat-message"
	BuyerPurchaseNotification    SourceEvent = "buyer-purchase"
	RemindToPayOrderNotification SourceEvent = "remind-to-pay-order"

	EmailType NotiType = "email"
	PushType  NotiType = "push"
)

var (
	DefaultNotiChannelBySourceEvent map[SourceEvent][]NotiType = map[SourceEvent][]NotiType{
		ItemShippedNotification:      {PushType},
		ChatMessageNotification:      {PushType, EmailType},
		BuyerPurchaseNotification:    {PushType, EmailType},
		RemindToPayOrderNotification: {PushType},
	}
)
