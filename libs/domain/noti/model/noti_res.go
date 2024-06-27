package model

type NotiResult struct {
	ID        string `json:"id"`
	ReqID     string `json:"reqId"`
	IsSuccess bool   `json:"isSuccess"`
	Reason    string `json:"reason"`
}
