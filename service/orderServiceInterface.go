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

	// q := domain.QueueDB{
	// 	Client: client,
	// }

	// tm := time.Now().Local()
	// ts := tm.Format("2006-01-02 15:04:05")

	// qData := domain.Queue{
	// 	TableName:    "order",
	// 	LockArgument: strconv.Itoa(req.Orders[0].AccountID),
	// 	Datetime:     ts,
	// }
	// err := q.Enqueue(qData)
	// if err != nil {
	// 	return nil, err.Error()
	// }

	// defer q.Dequeue(qData)

	for _, v := range req.Orders {

		orderIn := domain.OrderIn{
			AccountID:   req.AccountID,
			Timestamp:   v.Timestamp,
			Operation:   v.Operation,
			IssuerName:  v.IssuerName,
			TotalShares: v.TotalShares,
			SharePrice:  v.SharePrice,
		}

		result, err := o.repo.CreateOrder(&orderIn)

		if err != nil {
			currentBalance = result.Balance
			issuer = Issuer{
				IssuerName: result.IssuerName,
				BusinessErrors: []string{
					err.Error(),
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

func NewOrderServiceInterface(r domain.OrderRepository) *OrderServiceInterface {
	return &OrderServiceInterface{
		repo: r,
	}
}
