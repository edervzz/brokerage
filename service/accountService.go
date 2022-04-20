package service

import "brokerage/tech"

type AccountCreateRequest struct {
	Balance float32 `json:"cash"`
}

type AccountCreateResponse struct {
	AccountID int      `json:"id"`
	Cash      float32  `json:"cash"`
	Issuers   []string `json:"issuers"`
}

type AccountService interface {
	CreateAccount(*AccountCreateRequest) (*AccountCreateResponse, *tech.AppMess)
}
