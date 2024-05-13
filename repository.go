package main

import (
	"database/sql"
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
	return nil
}

func (p *PostgresRepository) GetAllAccounts() error {
	return nil
}

func (p *PostgresRepository) GetAccountById(id AccountID) error {
	return nil
}

func (p *PostgresRepository) UpdateAccount(a *Account) error {
	return nil
}

func (p *PostgresRepository) DeleteAccountById(id AccountID) error {
	return nil
}
