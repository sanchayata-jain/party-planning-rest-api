package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/config"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/tables"
	"github.com/go-chi/chi"
	"go.opencensus.io/plugin/ochttp"
)

type Server interface {
	ListenandServe() error
	Shutdown(ctx context.Context) error
}

type server struct {
	serv *http.Server
}

func New() Server {
	ctx := context.Background()
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	db, err := NewDatabase(ctx)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	err = Init(db)
	if err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}

	r := chi.NewRouter()

	repo := tables.NewRepository(db)
	service := tables.NewService(repo)
	ctrl := tables.NewController(service)

	r.MethodFunc(http.MethodGet, "/ping", handlerPing)
	r.MethodFunc(http.MethodGet, "/get_tables", ctrl.ListTables())
	r.MethodFunc(http.MethodPost, "/tables", ctrl.Create())

	return &server{
		serv: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.Port),
			Handler: r,
		},
	}

}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func (s *server) ListenandServe() error {
	log.Printf("connect to http://%s or http://localhost:%s/ping for ping", "127.0.0.1", "3000")
	return http.ListenAndServe(fmt.Sprintf("%s:%s", "127.0.0.1", "3000"), &ochttp.Handler{
		Handler: s.serv.Handler,
	})
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.serv.Shutdown(ctx)
}
