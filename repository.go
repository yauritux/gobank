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
	connStr := "user=postgres dbname=gobank sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database!")
		return nil, err
	}
	return &PostgresRepository{
		db: db,
	}, nil
}
