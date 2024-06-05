package usecases

import (
	"log"
	"os"

	"github.com/Dias221467/MicroServices/internal/domain/models"
	"github.com/Dias221467/MicroServices/internal/interfaces/adapters/postgres"
)

type BookUsecase struct {
	BookRepo *postgres.BookRepository
	logger   *log.Logger
}

func NewBookUsecase(bookRepo *postgres.BookRepository) *BookUsecase {
	return &BookUsecase{
		BookRepo: bookRepo,
		logger:   log.New(os.Stdout, "USECASE: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (u *BookUsecase) AddBook(book *models.Book) error {
	u.logger.Println("Adding book:", book)
	if err := u.BookRepo.AddBook(book); err != nil {
		u.logger.Println("Error adding book:", err)
		return err
	}
	u.logger.Println("Book added successfully:", book)
	return nil
}

func (u *BookUsecase) GetBooks() ([]*models.Book, error) {
	u.logger.Println("Retrieving books")
	books, err := u.BookRepo.GetBooks()
	if err != nil {
		u.logger.Println("Error retrieving books:", err)
		return nil, err
	}
	u.logger.Println("Books retrieved successfully")
	return books, nil
}

func (u *BookUsecase) GetBookByID(id int) (*models.Book, error) {
	u.logger.Println("Retrieving book by ID:", id)
	book, err := u.BookRepo.GetBookByID(id)
	if err != nil {
		u.logger.Println("Error retrieving book by ID:", err)
		return nil, err
	}
	u.logger.Println("Book retrieved successfully:", book)
	return book, nil
}

func (u *BookUsecase) UpdateBook(book *models.Book) error {
	u.logger.Println("Updating book:", book)
	if err := u.BookRepo.UpdateBook(book); err != nil {
		u.logger.Println("Error updating book:", err)
		return err
	}
	u.logger.Println("Book updated successfully:", book)
	return nil
}

func (u *BookUsecase) DeleteBook(id int) error {
	u.logger.Println("Deleting book by ID:", id)
	if err := u.BookRepo.DeleteBook(id); err != nil {
		u.logger.Println("Error deleting book:", err)
		return err
	}
	u.logger.Println("Book deleted successfully, ID:", id)
	return nil
}
