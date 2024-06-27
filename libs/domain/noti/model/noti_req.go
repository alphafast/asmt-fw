package model

type EmailPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Token string `json:"token"`
}

type NotiRequest struct {
	ID           string        `json:"id"`
	ReqID        string        `json:"reqId"`
	SourceEvent  SourceEvent   `json:"sourceEvent"`
	NotiType     NotiType      `json:"notiType"`
	EmailPayload *EmailPayload `json:"emailPayload,omitempty"`
	PushPayload  *PushPayload  `json:"pushPayload,omitempty"`
}
