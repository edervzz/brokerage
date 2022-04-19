package domain

import "github.com/jmoiron/sqlx"

type Migration struct {
	client *sqlx.DB
}

func (m *Migration) CreateTables() {
	_, err := m.client.Exec(`CREATE DATABASE IF NOT EXIST brokerage`)
	if err != nil {
		panic(err)
	}

	_, err = m.client.Exec(`CREATE TABLE brokerage.account (
		account_id INT auto_increment NOT NULL,
		balance DECIMAL NULL,
		last_timestamp varchar(12) NULL,
		CONSTRAINT account_PK PRIMARY KEY (account_id))`)
	if err != nil {
		panic(err)
	}

	_, err = m.client.Exec(`CREATE TABLE brokerage.issuer (
		account_id INT NOT NULL,
		issuer_name varchar(5) NOT NULL,
		qty INT NULL,
		CONSTRAINT issuer_PK PRIMARY KEY (account_id,issuer_name),
		CONSTRAINT issuer_FK FOREIGN KEY (account_id) REFERENCES brokerage.account(account_id))`)
	if err != nil {
		panic(err)
	}

	_, err = m.client.Exec(`CREATE TABLE brokerage.order (
		order_id INT auto_increment NOT NULL,
		account_id INT NOT NULL,		
		timestamp varchar(12) NOT NULL,
		operation varchar(5) NOT NULL,
		issuer_name varchar(5) NOT NULL,
		total_shares int NOT NULL,
		share_price decimal NOT NULL,
		CONSTRAINT order_PK PRIMARY KEY (order_id),
		CONSTRAINT order_UN UNIQUE KEY (order_id,account_id),
		CONSTRAINT order_FK FOREIGN KEY (account_id) REFERENCES brokerage.account(account_id))`)
	if err != nil {
		panic(err)
	}

}
