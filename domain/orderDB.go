package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
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

	// get balance
	currentBalance, lastOperTime, err := db.GetBalanceAndLastOperation(o.AccountID, o.IssuerName)
	if err != nil {
		return nil, err
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
	issuer, err := db.GetIssuer(o.AccountID, o.IssuerName)
	if err != nil {
		return baseOrder, err
	}

	// run checks
	err = db.ChecksOrder(*o, issuer, newBalance, lastOperTime)
	if err != nil {
		return baseOrder, err
	}

	result, err := db.MakeTX(o, issuer, float32(newBalance))
	if err != nil {
		fmt.Println(err)
		return baseOrder, err
	}

	return result, nil
}

func (db *OrderDB) MakeTX(o *OrderIn, issuers *issuer, newBalance float32) (*Order, error) {

	i, err := strconv.ParseInt(o.Timestamp, 10, 64)
	tm := time.Unix(i, 0)
	ts := tm.Local().Format("2006-01-02 15:04:05")

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
			SET qty=?,
			last_operation=?
			WHERE account_id=?
			AND issuer_name=?`,
			o.TotalShares, ts, o.AccountID, o.IssuerName)
		if err != nil {
			tx.Rollback()
			fmt.Println("error: CreateOrder: cannot insert issuer:", err.Error())
			return nil, errors.New("INVALID_OPERATION")
		}
	} else {
		_, err = tx.Exec(`INSERT INTO brokerage.issuer
			(account_id, issuer_name, last_operation, qty)
			VALUES(?,?,?,?)`,
			o.AccountID, o.IssuerName, ts, o.TotalShares)
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

func (db *OrderDB) GetBalanceAndLastOperation(acctId int, issuer string) (float32, *time.Time, error) {
	var currentBalance float32
	var lastOperation sql.NullString
	var lastOperTime time.Time

	row := db.client.QueryRow(`SELECT a.balance, i.last_operation  
		FROM brokerage.account a
		LEFT JOIN brokerage.issuer i
		ON a.account_id = i.account_id
		AND i.issuer_name = ?
		WHERE a.account_id = ?`,
		issuer, acctId,
	)

	err := row.Scan(&currentBalance, &lastOperation)
	if err != nil {
		err = errors.New("error: CreateOrder: impossible to get balance: " + err.Error())
		fmt.Println(err)
		return 0, nil, errors.New("INVALID_OPERATION")
	}

	// iTime, _ := strconv.ParseInt(lastOperation.String, 10, 64)
	// lastOperTime = time.Unix(iTime, 0)

	lastOperTime, err = time.ParseInLocation("2006-01-02 15:04:05", lastOperation.String, time.Local)

	return currentBalance, &lastOperTime, nil
}

func (db *OrderDB) ChecksOrder(o OrderIn, i *issuer, currentBalance float32, lastOperTime *time.Time) error {

	if BUY != o.Operation && SELL != o.Operation {
		return errors.New("INVALID_OPERATION")
	}

	if BUY == o.Operation && currentBalance <= float32(o.TotalShares)*float32(o.SharePrice) {
		return errors.New("INSUFFICIENT_FUNDS")
	}

	if SELL == o.Operation && i.qty < o.TotalShares {
		return errors.New("INSUFFICIENT_STOCKS")
	}

	iTime, _ := strconv.ParseInt(o.Timestamp, 10, 64)
	thisTime := time.Unix(iTime, 0)

	fmt.Println(thisTime.Format("2006-01-02 15:04:05"))
	fmt.Println(lastOperTime.Format("2006-01-02 15:04:05"))

	if lastOperTime.Unix() != 0 {
		d := lastOperTime.Unix() - thisTime.Unix()
		fmt.Println(d)
		if math.Abs(float64(d)) <= 300 {
			return errors.New("DUPLICATED_OPERATION")
		}
	}

	// if !(thisTime.Hour() >= 6 && thisTime.Minute() >= 0 && thisTime.Second() >= 0 &&
	// 	thisTime.Hour() <= 14 && thisTime.Minute() <= 59 && thisTime.Second() >= 59) {
	// 	return errors.New("CLOSET_MARKET")
	// }

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
		fmt.Println(err)
		return nil, errors.New("INVALID_OPERATION")
	}

	for rows.Next() {
		err = rows.Scan(&nameIssuer, &qtyIssuer)
		if err != nil {
			fmt.Println(err)
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

func NewOrderDB(c *sqlx.DB) *OrderDB {
	return &OrderDB{
		client: c,
	}
}
