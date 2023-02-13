package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(mysql:3306)/database?charset=utf8mb4&parseTime=True")
	if err != nil {
		log.Fatal("sql.Open: didn't create db")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		// log.Fatal("ping isn't pinging")
		log.Fatal(err)
	}

	// repo := tables.NewRepository(db)
	// service := tables.NewService(repo)
	// ctrl := tables.NewController(service)

	// ping
	http.HandleFunc("/ping", handlerPing)
	// http.HandleFunc("/tables", ctrl.Create())

	http.ListenAndServe(":3000", nil)
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
