package main

import (
	"encoding/json"
	"net/http"
)

type AccountService interface {
	Create(account Account) (Account, error)
	ReadAll() ([]Account, error)
	Read(id AccountID) (Account, error)
	Update(id AccountID, account Account) (Account, error)
	Delete(id AccountID) error
}

type AccountError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type APIServer struct {
	accountService AccountService
}

func NewAPIServer(as AccountService) *APIServer {
	return &APIServer{
		accountService: as,
	}
}

func (s *APIServer) errorResponse(w http.ResponseWriter, statusCode int, errorString string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encodingError := json.NewEncoder(w).Encode(AccountError{
		StatusCode: statusCode,
		Error:      errorString,
	})
	if encodingError != nil {
		http.Error(w, encodingError.Error(), http.StatusInternalServerError)
	}
}

func (s *APIServer) handleCreateAccounts(w http.ResponseWriter, r *http.Request) {
	var account Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		s.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newAccount, err := s.accountService.Create(account)
	if err != nil {
		s.errorResponse(w, http.StatusExpectationFailed, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newAccount)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.accountService.ReadAll()

	if err != nil {
		s.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	if r.URL.Query().Get("format") == "xlsx" {
		file, err := WriteAccountsToExcel(accounts)
		if err != nil {
			s.errorResponse(w, http.StatusExpectationFailed, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+"accounts.xlsx")
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		file.Write(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(accounts)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) {
	id := AccountID(r.PathValue("id"))
	account, err := s.accountService.Read(id)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := AccountID(r.PathValue("id"))

	var updatedAccount Account
	err := json.NewDecoder(r.Body).Decode(&updatedAccount)
	if err != nil {
		s.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := s.accountService.Update(id, updatedAccount)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleDeleteAccounts(w http.ResponseWriter, r *http.Request) {
	id := AccountID(r.PathValue("id"))
	err := s.accountService.Delete(id)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
