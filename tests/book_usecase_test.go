package tests

import (
	"database/sql"
	"testing"

	"github.com/Dias221467/MicroServices/internal/domain/models"
	"github.com/Dias221467/MicroServices/internal/interfaces/adapters/postgres"
	"github.com/Dias221467/MicroServices/internal/usecases"
	_ "github.com/lib/pq"
)

var db *sql.DB
var usecase *usecases.BookUsecase

func setup() {
	var err error
	db, err = sql.Open("postgres", "postgresql://postgres:lbfc2005@localhost:5432/microservices?sslmode=disable")
	if err != nil {
		panic(err)
	}

	bookRepo := postgres.NewBookRepository(db)
	usecase = usecases.NewBookUsecase(bookRepo)
}

func teardown() {
	db.Close()
}

func TestAddBook(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "Test Book",
		Author:   "Author Name",
		BookYear: 2022,
	}

	err := usecase.AddBook(book)
	if err != nil {
		t.Errorf("Failed to add book: %v", err)
	}

	if book.ID == 0 {
		t.Error("Expected non-zero book ID after adding book")
	}
}

func TestGetBooks(t *testing.T) {
	setup()
	defer teardown()

	books, err := usecase.GetBooks()
	if err != nil {
		t.Errorf("Failed to get books: %v", err)
	}

	if len(books) == 0 {
		t.Error("Expected at least one book")
	}
}

func TestGetBookByID(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "Test Book",
		Author:   "Author Name",
		BookYear: 2022,
	}
	usecase.AddBook(book)

	retrievedBook, err := usecase.GetBookByID(book.ID)
	if err != nil {
		t.Errorf("Failed to get book by ID: %v", err)
	}

	if retrievedBook == nil {
		t.Error("Expected non-nil book")
	}
	if retrievedBook.ID != book.ID {
		t.Errorf("Expected book ID %d, got %d", book.ID, retrievedBook.ID)
	}
}

func TestUpdateBook(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "Test Book",
		Author:   "Author Name",
		BookYear: 2022,
	}
	usecase.AddBook(book)

	book.Title = "Updated Title"
	err := usecase.UpdateBook(book)
	if err != nil {
		t.Errorf("Failed to update book: %v", err)
	}

	updatedBook, err := usecase.GetBookByID(book.ID)
	if err != nil {
		t.Errorf("Failed to get book by ID: %v", err)
	}
	if updatedBook.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", updatedBook.Title)
	}
}

func TestDeleteBook(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "Test Book",
		Author:   "Author Name",
		BookYear: 2022,
	}
	usecase.AddBook(book)

	err := usecase.DeleteBook(book.ID)
	if err != nil {
		t.Errorf("Failed to delete book: %v", err)
	}

	deletedBook, err := usecase.GetBookByID(book.ID)
	if deletedBook != nil {
		t.Error("Expected nil book after deletion")
	}
	if err == nil {
		t.Error("Expected error when retrieving deleted book")
	}
}

func TestAddBook_EmptyTitle(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "",
		Author:   "Author Name",
		BookYear: 2022,
	}

	err := usecase.AddBook(book)
	if err == nil {
		t.Error("Expected error when adding book with empty title")
	}
}

func TestAddBook_EmptyAuthor(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "Test Book",
		Author:   "",
		BookYear: 2022,
	}

	err := usecase.AddBook(book)
	if err == nil {
		t.Error("Expected error when adding book with empty author")
	}
}

func TestAddBook_InvalidYear(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		Title:    "Test Book",
		Author:   "Author Name",
		BookYear: -1, // Invalid year
	}

	err := usecase.AddBook(book)
	if err == nil {
		t.Error("Expected error when adding book with invalid year")
	}
}

func TestGetBookByID_NotFound(t *testing.T) {
	setup()
	defer teardown()

	_, err := usecase.GetBookByID(9999) // Assuming 9999 is a non-existent ID
	if err == nil {
		t.Error("Expected error when retrieving non-existent book")
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	setup()
	defer teardown()

	book := &models.Book{
		ID:       9999, // Assuming 9999 is a non-existent ID
		Title:    "Non-existent Book",
		Author:   "Author Name",
		BookYear: 2022,
	}

	err := usecase.UpdateBook(book)
	if err == nil {
		t.Error("Expected error when updating non-existent book")
	}
}
