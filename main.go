package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Accounts service at port 8080...")

	accountService := NewService()
	accountsHandler := NewAPIServer(accountService)

	err := http.ListenAndServe(":8080", initializeRoutes(accountsHandler))
	if err != nil {
		panic(err)
	}
}

func initializeRoutes(s *APIServer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/accounts", s.handleGetAccounts)
	mux.HandleFunc("GET /api/accounts/{id}", s.handleGetAccountById)
	mux.HandleFunc("POST /api/accounts", s.handleCreateAccounts)
	mux.HandleFunc("PUT /api/accounts/{id}", s.handleUpdateAccount)
	mux.HandleFunc("DELETE /api/accounts/{id}", s.handleDeleteAccounts)
	return mux
}
