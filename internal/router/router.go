package router

import (
	"database/sql"

	"github.com/Dias221467/MicroServices/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", handlers.GetBooks(db)).Methods("GET")
	r.HandleFunc("/books/{id}", handlers.GetBook(db)).Methods("GET")
	r.HandleFunc("/books", handlers.CreateBook(db)).Methods("POST")
	r.HandleFunc("/books/{id}", handlers.UpdateBook(db)).Methods("PUT")
	r.HandleFunc("/books/{id}", handlers.DeleteBook(db)).Methods("DELETE")
	return r
}
