package domain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccountDB struct {
	client *sqlx.DB
}

func (db AccountDB) CreateAccount(cash float32) (*Account, *BusinessError) {
	result, err := db.client.Exec(`INSERT INTO brokerage.account
		(balance)
		VALUES(?)`,
		cash)
	if err != nil {
		fmt.Println("error: CreateAccount: check if db was created")
		be := &BusinessError{
			Message: []string{
				"INVALID_OPERATION",
				"CHECK_LOGS",
			},
		}
		return nil, be
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("error: CreateAccount: impossible retrieve account id")
		be := &BusinessError{
			Message: []string{
				"INVALID_OPERATION",
				"CHECK_LOGS",
			},
		}
		return nil, be
	}

	return &Account{
		AccountID: int(id),
		Cash:      cash,
		Issuers:   []Issuer{},
	}, nil
}
