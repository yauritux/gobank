package main

import (
	"encoding/json"
	"net/http"
)

type AccountService interface {
	Create(account Account) Account
	ReadAll() []Account
	Read(id AccountID) (Account, bool)
	Update(id AccountID, account Account) (Account, bool)
	Delete(id AccountID) bool
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

	newAccount := s.accountService.Create(account)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newAccount)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts := s.accountService.ReadAll()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(accounts)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) {
	id := AccountID(r.PathValue("id"))
	account, found := s.accountService.Read(id)
	if !found {
		s.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(account)
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

	account, found := s.accountService.Update(id, updatedAccount)
	if !found {
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
	ok := s.accountService.Delete(id)
	if !ok {
		s.errorResponse(w, http.StatusNotFound, "Not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
