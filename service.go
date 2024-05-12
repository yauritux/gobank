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
	accounts map[AccountID]Account
}

func NewService() *Service {
	return &Service{
		accounts: make(map[AccountID]Account),
	}
}

func (s *Service) Create(account Account) Account {
	s.Lock()
	defer s.Unlock()

	id := NewAccountID()
	account.ID = string(id)
	s.accounts[id] = account
	return account
}

func (s *Service) ReadAll() []Account {
	s.Lock()
	defer s.Unlock()

	accounts := make([]Account, 0, len(s.accounts))
	for _, account := range s.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}

func (s *Service) Read(id AccountID) (Account, bool) {
	s.Lock()
	defer s.Unlock()

	account, ok := s.accounts[id]
	return account, ok
}

func (s *Service) Update(id AccountID, account Account) (Account, bool) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.accounts[id]
	if !ok {
		return Account{}, false
	}
	account.ID = string(id)
	s.accounts[id] = account
	return account, true
}

func (s *Service) Delete(id AccountID) bool {
	s.Lock()
	defer s.Unlock()

	_, ok := s.accounts[id]
	if !ok {
		return false
	}
	delete(s.accounts, id)
	return true
}
