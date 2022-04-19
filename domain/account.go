package domain

import "brokerage/tech"

type Issuer struct {
	Name string
	Qty  int
}

type Account struct {
	AccountID int
	Cash      float32
	Issuers   []Issuer
}

type AccountRepository interface {
	CreateAccount(cash float32) (*Account, *tech.AppMess)
}
