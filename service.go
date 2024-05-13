package main

import (
	"sync"

	"github.com/google/uuid"
)

type AccountID string

func NewAccountID() AccountID {
	return AccountID(uuid.NewString())
}

type Account struct {
	ID            string  `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

type Service struct {
	sync.Mutex
	db *PostgresRepository
}

func NewService(db *PostgresRepository) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Create(account Account) (Account, error) {
	s.Lock()
	defer s.Unlock()

	id := NewAccountID()
	account.ID = string(id)
	if err := s.db.CreateNewAccount(&account); err != nil {
		return account, err
	}
	return account, nil
}

func (s *Service) ReadAll() ([]Account, error) {
	s.Lock()
	defer s.Unlock()

	accounts, err := s.db.GetAllAccounts()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *Service) Read(id AccountID) (Account, error) {
	s.Lock()
	defer s.Unlock()

	return s.db.GetAccountById(id)
}

func (s *Service) Update(id AccountID, account Account) (Account, error) {
	s.Lock()
	defer s.Unlock()

	account.ID = string(id)
	err := s.db.UpdateAccount(&account)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (s *Service) Delete(id AccountID) error {
	s.Lock()
	defer s.Unlock()

	if err := s.db.DeleteAccountById(id); err != nil {
		return err
	}
	return nil
}
