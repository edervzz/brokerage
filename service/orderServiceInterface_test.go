package service

import (
	"brokerage/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Order(t *testing.T) {
	// given
	tt := []struct {
		name      string
		isError   bool
		AccountID int
		Orders    []OrderCreate
	}{
		{
			"ERR cannot create order",
			true,
			1,
			[]OrderCreate{
				{
					Timestamp:   "1650577527",
					Operation:   "SELL",
					IssuerName:  "APPL",
					TotalShares: 5,
					SharePrice:  10,
				},
			},
		},
		{
			"OK create order w/ 1 item",
			false,
			1,
			[]OrderCreate{
				{
					Timestamp:   "1650577527",
					Operation:   "SELL",
					IssuerName:  "APPL",
					TotalShares: 5,
					SharePrice:  10,
				},
			},
		},
		{
			"OK create order w/ many items",
			false,
			1,
			[]OrderCreate{
				{
					Timestamp:   "1650577527",
					Operation:   "SELL",
					IssuerName:  "APPL",
					TotalShares: 5,
					SharePrice:  10,
				},
				{
					Timestamp:   "1650577528",
					Operation:   "BUY",
					IssuerName:  "SBUX",
					TotalShares: 2,
					SharePrice:  8,
				},
			},
		},
	}

	for _, tc := range tt {
		service := NewOrderServiceInterface(domain.NewOrderStub(tc.isError))
		req := OrderCreateRequest{
			AccountID: tc.AccountID,
			Orders:    tc.Orders,
		}

		// when
		response := service.CreateOrder(&req)
		// then
		assert.GreaterOrEqual(t, response.CBalance.Cash, float32(0))
		for _, v := range response.CBalance.Issuers {
			assert.NotNil(t, v.BusinessErrors)
		}
	}

}
