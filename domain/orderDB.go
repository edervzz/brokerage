package domain

import (
	"brokerage/tech"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jmoiron/sqlx"
)

type issuer struct {
	name string
	qty  int
}

type OrderDB struct {
	client *sqlx.DB
}

const BUY string = "BUY"
const SELL string = "SELL"

func (db *OrderDB) CreateOrder(o *OrderIn) (*Order, error) {
	var currentBalance, newBalance float32

	if o.Timestamp == "" {
		o.Timestamp = fmt.Sprint(time.Now().Unix())
	}

	currentBalance, err := db.GetBalance(o.AccountID, o.IssuerName)
	if err != nil {
		return nil, err
	}

	newOrder := &Order{
		AccountID: o.AccountID,
		Balance:   float32(newBalance),
	}

	// set new balance
	switch o.Operation {
	case BUY:
		newBalance = currentBalance - (float32(o.TotalShares) * float32(o.SharePrice))
	case SELL:
		newBalance = currentBalance + (float32(o.TotalShares) * float32(o.SharePrice))
	}

	// get issuerInfo
	issuerInfo, err := db.GetIssuer(o.AccountID, o.IssuerName)
	if err != nil {
		return newOrder, err
	}
	// get last operation
	lastOperation, operation, _ := db.GetLastOperation(o.AccountID, o.IssuerName, o.TotalShares)

	// run checks
	err = db.ChecksOrder(*o, issuerInfo, newBalance, lastOperation, operation)
	if err != nil {
		return newOrder, err
	}

	result, err := db.MakeTX(o, issuerInfo, float32(newBalance))
	if err != nil {
		fmt.Println(err)
		return newOrder, err
	}

	return result, nil
}

func (db *OrderDB) MakeTX(o *OrderIn, issuers *issuer, newBalance float32) (*Order, error) {

	ts := tech.IntString2yymmdd_hhmmss(o.Timestamp)

	tx, err := db.client.Begin()

	_, err = tx.Exec(`INSERT INTO brokerage.order
		(account_id, timestamp, operation, issuer_name, total_shares, share_price)
		VALUES(?,?,?,?,?,?)`,
		o.AccountID, ts, o.Operation, o.IssuerName, o.TotalShares, o.SharePrice)
	if err != nil {
		tx.Rollback()
		fmt.Println("error: CreateOrder: cannot create order:", err.Error())
		return nil, errors.New("INVALID_OPERATION")
	}

	if issuers != nil && o.IssuerName == (*issuers).name {
		_, err = tx.Exec(`UPDATE brokerage.issuer
			SET qty=?
			WHERE account_id = ?
			AND issuer_name = ?`,
			o.TotalShares, o.AccountID, o.IssuerName)
		if err != nil {
			tx.Rollback()
			fmt.Println("error: CreateOrder: cannot update issuer:", err.Error())
			return nil, errors.New("INVALID_OPERATION")
		}
	} else {
		_, err = tx.Exec(`INSERT INTO brokerage.issuer
			(account_id, issuer_name, qty)
			VALUES(?,?,?)`,
			o.AccountID, o.IssuerName, o.TotalShares)
		if err != nil {
			tx.Rollback()
			fmt.Println("error: CreateOrder: cannot insert issuer:", err.Error())
			return nil, errors.New("INVALID_OPERATION")
		}
	}

	_, err = tx.Exec(`UPDATE brokerage.account
		SET balance = ?
		WHERE account_id = ?`,
		newBalance, o.AccountID)
	if err != nil {
		tx.Rollback()
		fmt.Println("error: CreateOrder: cannot update balance:", err.Error())
		return nil, errors.New("INVALID_OPERATION")
	}

	tx.Commit()

	return &Order{
		AccountID:   o.AccountID,
		Timestamp:   o.Timestamp,
		Operation:   o.Operation,
		IssuerName:  o.IssuerName,
		TotalShares: o.TotalShares,
		SharePrice:  o.SharePrice,
		Balance:     float32(newBalance),
	}, nil
}

func (db *OrderDB) GetBalance(acctId int, issuer string) (float32, error) {
	var currentBalance float32

	row := db.client.QueryRow(`SELECT balance 
		FROM brokerage.account
		WHERE account_id = ?`,
		acctId,
	)

	err := row.Scan(&currentBalance)
	if err != nil {
		err = errors.New("error: CreateOrder: impossible to get balance: " + err.Error())
		fmt.Println(err)
		return 0, errors.New("INVALID_OPERATION")
	}

	return currentBalance, nil
}

func (db *OrderDB) ChecksOrder(o OrderIn, i *issuer, newBalance float32,
	lastOperTime *time.Time, lastOperation string) error {

	if BUY != o.Operation && SELL != o.Operation {
		return errors.New("INVALID_OPERATION_TYPE")
	}

	if (BUY == o.Operation && newBalance <= float32(o.TotalShares)*float32(o.SharePrice)) ||
		newBalance < 0 {
		return errors.New("INSUFFICIENT_FUNDS")
	}

	if SELL == o.Operation && i.qty < o.TotalShares {
		return errors.New("INSUFFICIENT_STOCKS")
	}

	gmt, err := tech.IntString2Time(o.Timestamp)
	if err != nil {
		tech.LogInfo(err.Error())
		return errors.New("INVALID_OPERATION")
	}

	if lastOperTime != nil {
		if lastOperTime.Unix() != 0 {
			d := gmt.Unix() - lastOperTime.Unix()
			if math.Abs(float64(d)) <= 300 &&
				o.Operation == lastOperation &&
				o.TotalShares == o.TotalShares {
				return errors.New("DUPLICATED_OPERATION")
			}
		}
	}

	localTime := gmt.Local()
	if !(localTime.Hour() >= 0 && localTime.Minute() >= 0 && localTime.Second() >= 0 &&
		localTime.Hour() <= 23 && localTime.Minute() <= 59 && localTime.Second() <= 59) {
		return errors.New("CLOSET_MARKET")
	}
	return nil
}

func (db *OrderDB) GetIssuer(acctId int, issuer_name string) (*issuer, error) {
	var nameIssuer string = ""
	var qtyIssuer int = 0

	rows, err := db.client.Query(`SELECT issuer_name, qty 
		FROM brokerage.issuer
		WHERE account_id=?`,
		acctId,
	)
	if err != nil {
		tech.LogInfo(err.Error())
		return nil, errors.New("INVALID_OPERATION")
	}

	for rows.Next() {
		err = rows.Scan(&nameIssuer, &qtyIssuer)
		if err != nil {
			tech.LogInfo(err.Error())
			return nil, errors.New("INVALID_OPERATION")
		}

		if nameIssuer == issuer_name {
			i := &issuer{
				name: nameIssuer,
				qty:  qtyIssuer,
			}
			return i, nil
		}
	}
	return nil, nil
}

func (db *OrderDB) GetLastOperation(accId int, issuer string, totalShares int) (*time.Time, string, error) {
	row := db.client.QueryRow(`SELECT MAX(timestamp) as maxts, operation 
		FROM brokerage.order
		WHERE account_id = ?
		AND issuer_name = ?
		AND total_shares = ?
		GROUP BY operation, total_shares
		ORDER BY maxts DESC`,
		accId, issuer, totalShares,
	)
	var timestamp, operation sql.NullString
	err := row.Scan(&timestamp, &operation)
	if err != nil {
		tech.LogInfo(err.Error())
		return nil, "", err
	}

	if timestamp.String == "" || operation.String == "" {
		return nil, "", nil
	}

	t, err := tech.String2Time(timestamp.String)
	if err != nil {
		tech.LogInfo(err.Error())
		return nil, "", err
	}

	return t, operation.String, nil
}

func NewOrderDB(c *sqlx.DB) *OrderDB {
	return &OrderDB{
		client: c,
	}
}
