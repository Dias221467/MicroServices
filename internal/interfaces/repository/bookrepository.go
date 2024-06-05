package repository

import "github.com/Dias221467/MicroServices/internal/domain/models"

type BookRepository interface {
	Create(book *models.Book) error
	FindAll() ([]*models.Book, error)
	FindByID(id int) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id int) error
}
