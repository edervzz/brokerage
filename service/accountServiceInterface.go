package service

import (
	"brokerage/domain"
	"brokerage/tech"
)

type AccountServiceInterface struct {
	repo domain.AccountRepository
}

func (a *AccountServiceInterface) CreateAccount(req AccountCreateRequest) (*AccountCreateResponse, *tech.AppMess) {

	if req.balance <= 0 {
		return nil, &tech.AppMess{
			Code:    400,
			Message: "balance must be greater than zero",
		}
	}

	acct, err := a.repo.CreateAccount(req.balance)
	if err != nil {
		return nil, &tech.AppMess{
			Code:    500,
			Message: "cannot create account",
		}
	}

	return &AccountCreateResponse{
		AccountID: acct.AccountID,
		Cash:      acct.Cash,
		Issuers:   []string{},
	}, nil

}
