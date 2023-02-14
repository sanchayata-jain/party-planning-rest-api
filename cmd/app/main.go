package main

import (
	// "database/sql"
	"context"
	"fmt"
	"log"

	// "log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/getground/tech-tasks/backend/cmd/app/internal/models"
	"github.com/getground/tech-tasks/backend/cmd/app/internal/tables"
)

func main() {
	// db, err := sql.Open("mysql", "user:password@tcp(mysql:3306)/database?charset=utf8mb4&parseTime=True")
	// if err != nil {
	// 	log.Fatal("sql.Open: didn't create db")
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	// log.Fatal("ping isn't pinging")
	// 	log.Fatal(err)
	// }
	ctx := context.Background()
	// r := chi.NewRouter()
	db, err := NewDatabase(ctx)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	err = Init(db)
	if err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}

	repo := tables.NewRepository(db)
	service := tables.NewService(repo)
	ctrl := tables.NewController(service)

	// ping
	http.HandleFunc("/ping", handlerPing)
	http.HandleFunc("/get_tables", ctrl.ListTables())
	// r.Method(http.MethodPost, "/tables", ctrl.Create())
	http.HandleFunc("/tables", ctrl.Create())

	http.ListenAndServe(":3000", nil)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func NewDatabase(ctx context.Context) (*gorm.DB, error) {
	psqlconn := "postgresql://user:password@localhost:5432/database?sslmode=disable"
	return gorm.Open(postgres.New(postgres.Config{
		DSN: psqlconn,
	}), &gorm.Config{})
}

func Init(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Table{},
	)
}
