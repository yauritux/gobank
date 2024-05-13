package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgres() (*PostgresRepository, error) {
	connStr := "user=gobank dbname=gobank host=local_postgres password=tux123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Failed to connect to database!")
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (p *PostgresRepository) CreateNewAccount(a *Account) error {
	insertStatement := `INSERT INTO accounts(id, first_name, last_name, account_number, balance)
	VALUES($1, $2, $3, $4, $5);`
	_, err := p.db.Exec(insertStatement, a.ID, a.FirstName, a.LastName, a.AccountNumber, a.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) GetAllAccounts() ([]Account, error) {
	rows, err := p.db.Query(`SELECT * FROM accounts;`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accounts []Account

	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.AccountNumber, &a.Balance); err != nil {
			return accounts, err
		}
		accounts = append(accounts, a)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}

	return accounts, nil
}

func (p *PostgresRepository) GetAccountById(id AccountID) (Account, error) {
	sqlStatement := `SELECT * FROM accounts WHERE id=$1;`
	var a Account
	var row = p.db.QueryRow(sqlStatement, string(id))
	err := row.Scan(&a.ID, &a.FirstName, &a.LastName, &a.AccountNumber, &a.Balance)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return a, err
	case nil:
		fmt.Println(a)
		return a, nil
	default:
		return a, err
	}
}

func (p *PostgresRepository) UpdateAccount(a *Account) error {
	sqlStatement := `UPDATE accounts SET first_name=$2, last_name=$3, account_number=$4, balance=$5
	WHERE id=$1;`
	_, err := p.db.Exec(sqlStatement, a.ID, a.FirstName, a.LastName, a.AccountNumber, a.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) DeleteAccountById(id AccountID) error {
	sqlStatement := `DELETE FROM accounts WHERE id=$1`
	_, err := p.db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}
