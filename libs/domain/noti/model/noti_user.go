package model

type NotiUserEmailChannelPayload struct {
	EmailAddress string `json:"emailAddress"`
}

type NotiUserPushChannelPayload struct {
	Token string `json:"token"`
}

type NotiUserNotiChannel struct {
	NotiType            NotiType                     `json:"notiType"`
	EmailChannelPayload *NotiUserEmailChannelPayload `json:"emailChannelPayload,omitempty"`
	PushChannelPayload  *NotiUserPushChannelPayload  `json:"pushChannelPayload,omitempty"`
}

type NotiUser struct {
	UserID   string                `json:"userId"`
	Channels []NotiUserNotiChannel `json:"channels"`
}

func (u *NotiUser) GetTypeChannelMap() map[NotiType]NotiUserNotiChannel {
	channels := make(map[NotiType]NotiUserNotiChannel)
	for _, c := range u.Channels {
		channels[c.NotiType] = c
	}
	return channels
}
