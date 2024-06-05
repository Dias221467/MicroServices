package interfaces

import (
	"net/http"

	"github.com/Dias221467/MicroServices/internal/domain/models"
)

// BookRepository defines the methods that any type of book repository must implement.
type BookRepository interface {
	Create(book *models.Book) error
	FindAll() ([]*models.Book, error)
	FindByID(id int) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id int) error
}

// BookUsecase defines the methods that any type of book usecase must implement.
type BookUsecase interface {
	AddBook(book *models.Book) error
	GetBooks() ([]*models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}

// BookHandler defines the methods that any type of book handler must implement.
type BookHandler interface {
	CreateBook(w http.ResponseWriter, r *http.Request)
	GetBooks(w http.ResponseWriter, r *http.Request)
	GetBook(w http.ResponseWriter, r *http.Request)
	UpdateBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
}
