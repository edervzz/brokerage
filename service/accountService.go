package service

import "brokerage/tech"

type AccountCreateRequest struct {
	balance float32
}

type AccountCreateResponse struct {
	AccountID int
	Cash      float32
	Issuers   []string
}

type AccountService interface {
	CreateAccount(*AccountCreateRequest) (*AccountCreateResponse, *tech.AppMess)
}
