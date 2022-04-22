package domain

import (
	"errors"
	"fmt"
	"time"
)

type OrderStubOk struct {
	isErr bool
}

func (db *OrderStubOk) CreateOrder(o *OrderIn) (*Order, error) {
	tmst := o.Timestamp
	if tmst == "" {
		tmst = fmt.Sprint(time.Now().Unix())
	}

	if db.isErr {
		return &Order{
			AccountID:     o.AccountID,
			OrderID:       5001,
			Timestamp:     o.Timestamp,
			Operation:     o.Operation,
			IssuerName:    o.IssuerName,
			TotalShares:   o.TotalShares,
			SharePrice:    o.SharePrice,
			Balance:       4900,
			BusinessError: []string{},
		}, errors.New("")
	}

	return &Order{
		AccountID:     o.AccountID,
		OrderID:       5001,
		Timestamp:     o.Timestamp,
		Operation:     o.Operation,
		IssuerName:    o.IssuerName,
		TotalShares:   o.TotalShares,
		SharePrice:    o.SharePrice,
		Balance:       4900,
		BusinessError: []string{},
	}, nil
}

func NewOrderStub(isErr bool) *OrderStubOk {
	return &OrderStubOk{
		isErr,
	}
}
