package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	fmt.Printf("Starting Accounts service at port %s...\n", os.Getenv("APP_CONTAINER_PORT"))

	tls, tlsEnabled := os.LookupEnv("APP_TLS")

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

	mux := initializeRoutes(accountsHandler)

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ts := time.Now().UTC().Format(time.RFC3339)

		if err := db.db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status": "error",
				"timestamp": ts,
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"timestamp": ts,
		})
	})

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_CONTAINER_PORT")), mux)
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
