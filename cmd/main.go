package main

import (
	"database/sql"

	"log"
	"net/http"

	"github.com/Dias221467/MicroServices/internal/router"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:lbfc2005@localhost:5432/microservices?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := router.SetupRouter(db)
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
