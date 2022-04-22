package service

import (
	"brokerage/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AccountCreation(t *testing.T) {
	// given
	tt := []struct {
		name  string
		isErr bool
		cash  float32
	}{
		{
			"ERR Account creation",
			true,
			0,
		},
		{
			"OK Account creation",
			false,
			5000,
		},
	}

	for _, tc := range tt {
		service := NewAccountServiceInterface(domain.NewAccountStubOk())
		req := &AccountCreateRequest{
			Balance: tc.cash,
		}
		// when
		res, appMess := service.CreateAccount(req)

		// then
		if tc.isErr {
			assert.Nil(t, res)
			assert.NotNil(t, appMess)
		} else {
			assert.NotNil(t, res)
			assert.Nil(t, appMess)
		}
	}
}

func Test_AccountCreationDBDown(t *testing.T) {
	// given
	tt := []struct {
		name  string
		isErr bool
		cash  float32
	}{
		{
			"ERR Account creation",
			true,
			5000,
		},
	}

	for _, tc := range tt {
		service := NewAccountServiceInterface(domain.NewAccountStubErr())
		req := &AccountCreateRequest{
			Balance: tc.cash,
		}
		// when
		res, appMess := service.CreateAccount(req)

		// then
		if tc.isErr {
			assert.Nil(t, res)
			assert.NotNil(t, appMess)
		} else {
			assert.NotNil(t, res)
			assert.Nil(t, appMess)
		}
	}
}
