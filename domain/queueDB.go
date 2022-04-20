package domain

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Queue struct {
	TableName    string
	LockArgument string
	Datetime     string
}

type QueueDB struct {
	Client *sqlx.DB
}

func (db *QueueDB) Enqueue(q Queue) error {
	result, err := db.Client.Exec(`INSERT INTO brokerage.queue
		(table_name, lock_argument, datetime)
		VALUES(?,?,?)`,
		q.TableName, q.LockArgument, q.Datetime,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows != 1 {
		return errors.New("error: cannot lock: " + q.TableName + " | " + q.LockArgument)
	}
	return nil
}

func (db *QueueDB) Dequeue(q Queue) error {
	result, err := db.Client.Exec(`DELETE FROM brokerage.queue
		WHERE table_name = ?
		AND lock_argument = ?`,
		q.TableName, q.LockArgument,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows != 1 {
		return errors.New("error: cannot unlock: " + q.TableName + " | " + q.LockArgument)
	}
	return nil
}
