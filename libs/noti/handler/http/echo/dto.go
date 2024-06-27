package echo

type NotiResultItemResponse struct {
	ID        string `json:"id"`
	IsSuccess bool   `json:"isSuccess"`
	Reason    string `json:"reason"`
}

type NotiResultResponse struct {
	ReqID string                   `json:"reqId"`
	Items []NotiResultItemResponse `json:"items"`
}
