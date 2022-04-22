package domain

import (
	"brokerage/tech"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	client *sqlx.DB
}

type MigrationDB interface {
	CreateTables()
}

func (m *Migration) CreateTables() error {
	var brokerage string

	row := m.client.QueryRow(`SELECT schema_name
	FROM information_schema.schemata
	WHERE schema_name LIKE 'brokerage'`)

	row.Scan(&brokerage)

	if brokerage != "brokerage" {

		_, err := m.client.Exec(`CREATE DATABASE brokerage`)
		if err != nil {
			tech.LogWarn(err.Error())
			return err
		}

		_, err = m.client.Exec(`CREATE TABLE brokerage.account (
			account_id INT auto_increment NOT NULL,
			balance DECIMAL NULL,
			CONSTRAINT account_PK PRIMARY KEY (account_id))`)
		if err != nil {
			tech.LogWarn(err.Error())
			return err
		}

		_, err = m.client.Exec(`CREATE TABLE brokerage.issuer (
			account_id INT NOT NULL,
			issuer_name varchar(5) NOT NULL,
			qty INT NULL,
			CONSTRAINT issuer_PK PRIMARY KEY (account_id,issuer_name),
			CONSTRAINT issuer_FK FOREIGN KEY (account_id) REFERENCES brokerage.account(account_id))`)
		if err != nil {
			tech.LogWarn(err.Error())
			return err
		}

		_, err = m.client.Exec(`CREATE TABLE brokerage.order (
			order_id INT auto_increment NOT NULL,
			account_id INT NOT NULL,		
			timestamp timestamp NOT NULL,
			operation varchar(5) NOT NULL,
			issuer_name varchar(5) NOT NULL,
			total_shares int NOT NULL,
			share_price decimal NOT NULL,
			CONSTRAINT order_PK PRIMARY KEY (order_id),
			CONSTRAINT order_UN UNIQUE KEY (order_id,account_id),
			CONSTRAINT order_FK FOREIGN KEY (account_id) REFERENCES brokerage.account(account_id))`)
		if err != nil {
			tech.LogWarn(err.Error())
			return err
		}

		_, err = m.client.Exec(`CREATE TABLE brokerage.queue (
			table_name varchar(50) NOT NULL,
			lock_argument varchar(100) NOT NULL,
			datetime varchar(12) NOT NULL,
			CONSTRAINT order_PK PRIMARY KEY (table_name,lock_argument))`)
		if err != nil {
			panic(err)
		}

		_, err = m.client.Exec(`CREATE INDEX order_account_id_IDX 
			USING BTREE 
			ON brokerage.order
			(account_id,issuer_name,total_shares)`)
		if err != nil {
			tech.LogWarn(err.Error())
			return err
		}

	}
	tech.LogInfo("info: database brokerage created")
	return nil
}

func NewMigration(client *sqlx.DB) *Migration {
	return &Migration{
		client,
	}
}
