package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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

func WriteAccountsToExcel(accounts []Account) (*excelize.File, error) {
	file := excelize.NewFile()

	headers := []string{"ID", "First Name", "Last Name", "Account Number", "Balance"}
	for i, header := range headers {
		file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", string(rune(65+i)), 1), header)
	}

	for i, row := range accounts {
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), row.ID)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), row.FirstName)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), row.LastName)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), row.AccountNumber)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), row.Balance)
	}

	err := file.SaveAs("accounts.xlsx")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return file, nil
}
