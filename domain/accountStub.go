package domain

import "errors"

type AccountStubOk struct{}

func (db AccountStubOk) CreateAccount(cash float32) (*Account, error) {
	return &Account{
		Cash:      cash,
		AccountID: 9999,
		Issuers:   []Issuer{},
	}, nil
}

func NewAccountStubOk() *AccountStubOk {
	return &AccountStubOk{}
}

// *************************************
type AccountStubErr struct{}

func (db AccountStubErr) CreateAccount(cash float32) (*Account, error) {
	return nil, errors.New("INVALID_OPERATION")
}

func NewAccountStubErr() *AccountStubErr {
	return &AccountStubErr{}
}
