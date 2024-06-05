package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Dias221467/MicroServices/internal/domain/models"
	adapters "github.com/Dias221467/MicroServices/internal/interfaces/adapters/postgres"
	"github.com/Dias221467/MicroServices/internal/usecases"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Importing the pq driver for PostgreSQL
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *httptest.Server {
	r := mux.NewRouter()
	db, err := sql.Open("postgres", "postgresql://postgres:lbfc2005@localhost:5432/microservices?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	bookRepo := adapters.NewBookRepository(db)
	bookUsecase := usecases.NewBookUsecase(bookRepo)

	r.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var book models.Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = bookUsecase.AddBook(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(book)
		} else if r.Method == http.MethodGet {
			books, err := bookUsecase.GetBooks()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(books)
		}
	}).Methods(http.MethodPost, http.MethodGet)

	r.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			book, err := bookUsecase.GetBookByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(book)
		} else if r.Method == http.MethodPut {
			var book models.Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			book.ID = id
			err = bookUsecase.UpdateBook(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(book)
		} else if r.Method == http.MethodDelete {
			err := bookUsecase.DeleteBook(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		}
	}).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)

	return httptest.NewServer(r)
}

func TestIntegration_CreateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	book := &models.Book{
		Title:    "Integration Test Book",
		Author:   "Test Author",
		BookYear: 2024,
	}

	data, _ := json.Marshal(book)
	resp, err := http.Post(server.URL+"/books", "application/json", bytes.NewBuffer(data))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdBook models.Book
	json.NewDecoder(resp.Body).Decode(&createdBook)
	assert.NotZero(t, createdBook.ID)
	assert.Equal(t, book.Title, createdBook.Title)
	assert.Equal(t, book.Author, createdBook.Author)
	assert.Equal(t, book.BookYear, createdBook.BookYear)
}

func TestIntegration_GetBooks(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/books")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var books []models.Book
	json.NewDecoder(resp.Body).Decode(&books)
	assert.GreaterOrEqual(t, len(books), 1)
}

func TestIntegration_GetBookByID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First, create a new book
	book := &models.Book{
		Title:    "Integration Test Book",
		Author:   "Test Author",
		BookYear: 2024,
	}

	data, _ := json.Marshal(book)
	createResp, _ := http.Post(server.URL+"/books", "application/json", bytes.NewBuffer(data))

	var createdBook models.Book
	json.NewDecoder(createResp.Body).Decode(&createdBook)

	// Now, retrieve the book by ID
	resp, err := http.Get(server.URL + "/books/" + strconv.Itoa(createdBook.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedBook models.Book
	json.NewDecoder(resp.Body).Decode(&retrievedBook)
	assert.Equal(t, createdBook.ID, retrievedBook.ID)
	assert.Equal(t, createdBook.Title, retrievedBook.Title)
	assert.Equal(t, createdBook.Author, retrievedBook.Author)
	assert.Equal(t, createdBook.BookYear, retrievedBook.BookYear)
}

func TestIntegration_UpdateBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First, create a new book
	book := &models.Book{
		Title:    "Integration Test Book",
		Author:   "Test Author",
		BookYear: 2024,
	}

	data, _ := json.Marshal(book)
	createResp, _ := http.Post(server.URL+"/books", "application/json", bytes.NewBuffer(data))

	var createdBook models.Book
	json.NewDecoder(createResp.Body).Decode(&createdBook)

	// Now, update the book
	createdBook.Title = "Updated Title"
	updateData, _ := json.Marshal(createdBook)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/books/"+strconv.Itoa(createdBook.ID), bytes.NewBuffer(updateData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedBook models.Book
	json.NewDecoder(resp.Body).Decode(&updatedBook)
	assert.Equal(t, "Updated Title", updatedBook.Title)
}

func TestIntegration_DeleteBook(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// First, create a new book
	book := &models.Book{
		Title:    "Integration Test Book",
		Author:   "Test Author",
		BookYear: 2024,
	}

	data, _ := json.Marshal(book)
	createResp, _ := http.Post(server.URL+"/books", "application/json", bytes.NewBuffer(data))

	var createdBook models.Book
	json.NewDecoder(createResp.Body).Decode(&createdBook)

	// Now, delete the book
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/books/"+strconv.Itoa(createdBook.ID), nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify the book has been deleted
	getResp, _ := http.Get(server.URL + "/books/" + strconv.Itoa(createdBook.ID))
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}
