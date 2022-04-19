package service

import (
	"brokerage/domain"
	"sort"
)

type OrderServiceInterface struct {
	repo domain.OrderRepository
}

func (o *OrderServiceInterface) CreateOrder(req *OrderCreateRequest) *OrderCreateResponse {
	var currentBalance float32
	var issuer Issuer
	var issuers []Issuer

	// first sell to retrieve more funds before buy
	sort.SliceStable(req.Orders, func(i, j int) bool {
		return req.Orders[i].Operation > req.Orders[j].Operation
	})

	for _, v := range req.Orders {

		orderIn := domain.OrderIn{
			AccountID:   v.AccountID,
			Timestamp:   v.Timestamp,
			Operation:   v.Operation,
			IssuerName:  v.IssuerName,
			TotalShares: v.TotalShares,
			SharePrice:  v.SharePrice,
		}

		result, be := o.repo.CreateOrder(orderIn)

		if be != "" {
			currentBalance = result.Balance
			issuer = Issuer{
				IssuerName: result.IssuerName,
				BusinessErrors: []string{
					be,
				},
			}
			issuers = append(issuers, issuer)
		} else {
			currentBalance = result.Balance
			issuer = Issuer{
				IssuerName:     result.IssuerName,
				TotalShares:    result.TotalShares,
				SharePrice:     result.SharePrice,
				BusinessErrors: []string{},
			}
			issuers = append(issuers, issuer)
		}
	}

	return &OrderCreateResponse{
		CBalance: CurrentBalance{
			Cash:    currentBalance,
			Issuers: issuers,
		},
	}

}
