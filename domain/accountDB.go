package domain

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccountDB struct {
	client *sqlx.DB
}

func (db AccountDB) CreateAccount(cash float32) (*Account, error) {
	result, err := db.client.Exec(`INSERT INTO brokerage.account
		(balance)
		VALUES(?)`,
		cash)
	if err != nil {
		fmt.Println("error: CreateAccount: check if db was created:" + err.Error())
		return nil, errors.New("INVALID_OPERATION")
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error: CreateAccount: impossible retrieve account id:" + err.Error())
		return nil, errors.New("INVALID_OPERATION")
	}

	return &Account{
		AccountID: int(id),
		Cash:      cash,
		Issuers:   []Issuer{},
	}, nil
}

func NewAccountDB(client *sqlx.DB) *AccountDB {
	return &AccountDB{
		client,
	}
}
