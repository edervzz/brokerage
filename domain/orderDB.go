package domain

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type issuers struct {
	name string
	qty  int
}

type OrderDB struct {
	client *sqlx.DB
}

const BUY string = "BUY"
const SELL string = "SELL"

func (db *OrderDB) CreateOrder(o OrderIn) (*Order, string) {
	var currentBalance, newBalance float32

	// get balance
	currentBalance, lastTimestamp, be := db.GetBalanceAndLastTime(o.AccountID)
	if be != "" {
		return nil, be
	}

	baseOrder := &Order{
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

	// get issuer
	issuer, be := db.GetIssuer(o.AccountID, o.IssuerName)
	if be != "" {
		return baseOrder, be
	}

	// run checks
	be = db.ChecksOrder(o, issuer, newBalance, lastTimestamp)
	if be != "" {
		return baseOrder, be
	}

	result, be := db.MakeTX(o, issuer, float32(newBalance))
	if be != "" {
		fmt.Println(be)
		return baseOrder, be
	}

	return result, ""
}

func (db *OrderDB) MakeTX(o OrderIn, issuers *issuers, newBalance float32) (*Order, string) {

	tx, err := db.client.Begin()

	_, err = tx.Exec(`INSERT INTO brokerage.order
		(account_id, timestamp, operation, issuer_name, total_shares, share_price)
		VALUES(?,?,?,?,?)`,
		o.AccountID, o.Timestamp, o.IssuerName, o.TotalShares, o.SharePrice)
	if err != nil {
		tx.Rollback()
		fmt.Println("error: CreateOrder: cannot create order")
		return nil, "INVALID_OPERATION, CHECK_LOGS"
	}

	if o.IssuerName == (*issuers).name {
		_, err = tx.Exec(`UPDATE brokerage.issuers
			SET qty=?
			WHERE account_id=?
			AND issuer_name=?`,
			o.TotalShares, o.AccountID, o.IssuerName)
		if err != nil {
			tx.Rollback()
			fmt.Println("error: CreateOrder: cannot insert issuer")
			return nil, "INVALID_OPERATION, CHECK_LOGS"
		}
	} else {
		_, err = tx.Exec(`INSERT INTO brokerage.issuers
			(account_id, issuer_name, qty)
			VALUES(?,?,?)`,
			o.AccountID, o.IssuerName, o.TotalShares)
		if err != nil {
			tx.Rollback()
			fmt.Println("error: CreateOrder: cannot insert issuer")
			return nil, "INVALID_OPERATION, CHECK_LOGS"
		}
	}

	_, err = tx.Exec(`UPDATE brokerage.account
		SET balance = ?
			last_timestamp = ?
		WHERE account_id = ?`,
		newBalance, o.Timestamp, o.AccountID)
	if err != nil {
		tx.Rollback()
		fmt.Println("error: CreateOrder: cannot update balance")
		return nil, "INVALID_OPERATION, CHECK_LOGS"
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
	}, ""
}

func (db *OrderDB) GetBalanceAndLastTime(acctId int) (float32, string, string) {
	var currentBalance float32
	var lastTimestamp string

	row := db.client.QueryRow(`SELECT balance, last_timestamp
		FROM brokerage.account
		WHERE account_id=?`,
		acctId,
	)

	err := row.Scan(&currentBalance, &lastTimestamp)
	if err != nil {
		err = errors.New("error: CreateOrder: impossible to get balance")
		fmt.Println(err)
		return 0, "", "INVALID_OPERATION, CHECK_LOGS"
	}
	return currentBalance, lastTimestamp, ""
}

func (db *OrderDB) ChecksOrder(o OrderIn, i *issuers, currentBalance float32, lastTimestamp string) string {

	if BUY != o.Operation && SELL != o.Operation {
		return "INVALID_OPERATION"
	}

	if BUY == o.Operation && currentBalance <= float32(o.TotalShares)*float32(o.SharePrice) {
		return "INSUFFICIENT_FUNDS"
	}

	if SELL == o.Operation && i.qty < o.TotalShares {
		return "INSUFFICIENT_STOCKS"
	}

	iTime, _ := strconv.ParseInt(lastTimestamp, 10, 64)
	lastTime := time.Unix(iTime, 0)

	iTime, _ = strconv.ParseInt(o.Timestamp, 10, 64)
	thisTime := time.Unix(iTime, 0)

	duration := thisTime.Sub(lastTime)

	if !(duration > 5) {

		return "DUPLICATED_OPERATION"
	}

	if !(thisTime.Hour() >= 6 && thisTime.Minute() >= 0 && thisTime.Second() >= 0 &&
		thisTime.Hour() <= 14 && thisTime.Minute() <= 59 && thisTime.Second() >= 59) {
		return "CLOSET_MARKET"
	}

	return ""
}

func (db *OrderDB) GetIssuer(acctId int, issuer_name string) (*issuers, string) {
	var nameIssuer string = ""
	var qtyIssuer int = 0

	rows, err := db.client.Query(`SELECT name, qty 
		FROM issuers
		WHERE account_id=?`,
		acctId,
	)
	if err != nil {
		fmt.Println(err)
		return nil, "INVALID_OPERATION, CHECK_LOGS"
	}

	for rows.Next() {
		err = rows.Scan(&nameIssuer, &qtyIssuer)
		if err != nil {
			fmt.Println(err)
			return nil, "INVALID_OPERATION, CHECK_LOGS"
		}

		if nameIssuer == issuer_name {
			i := &issuers{
				name: nameIssuer,
				qty:  qtyIssuer,
			}
			return i, ""
		}
	}
	return nil, ""
}
