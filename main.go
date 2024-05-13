package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type OptFunc func(*Opts)

type Opts struct {
	maxConn int
	id      string
	tls     bool
}

func defaultOpts() Opts {
	return Opts{
		maxConn: 10,
		id:      "gobank",
		tls:     false,
	}
}

func WithTLS(opts *Opts) {
	opts.tls = true
}

func WithMaxConn(n int) OptFunc {
	return func(opts *Opts) {
		opts.maxConn = n
	}
}

func WithID(id string) OptFunc {
	return func(opts *Opts) {
		opts.id = id
	}
}

type Server struct {
	Opts
}

func NewServer(opts ...OptFunc) *Server {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}
	return &Server{
		Opts: o,
	}
}

func (s *Server) SetServerID(id string) {
	s.id = id
}

func (s *Server) SetMaxConn(n int) {
	s.maxConn = n
}

func main() {
	fmt.Println("Starting Accounts service at port 8080...")

	tls, tlsEnabled := os.LookupEnv("TLS")

	s := NewServer()

	if tlsEnabled && tls == "1" {
		s = NewServer(WithTLS, WithID("gobank-svc"))
	}

	maxConnStr, ok := os.LookupEnv("MAX_CONN")
	if ok {
		maxConn, err := strconv.Atoi(maxConnStr)
		if err != nil {
			panic(err)
		}
		s.SetMaxConn(maxConn)
	}

	fmt.Printf("%+v\n", s)

	db, err := NewPostgres()
	if err != nil {
		log.Fatal(err)
		return
	}

	accountService := NewService(db)
	accountsHandler := NewAPIServer(accountService)

	err = http.ListenAndServe(":8080", initializeRoutes(accountsHandler))
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
