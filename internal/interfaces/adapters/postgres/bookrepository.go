package postgres

import (
	"database/sql"

	"github.com/Dias221467/MicroServices/internal/domain/models"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) AddBook(book *models.Book) error {
	query := `INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING id`
	return r.DB.QueryRow(query, book.Title, book.Author, book.BookYear).Scan(&book.ID)
}

func (r *BookRepository) GetBooks() ([]*models.Book, error) {
	rows, err := r.DB.Query(`SELECT id, title, author, year FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.BookYear); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	return books, nil
}

func (r *BookRepository) GetBookByID(id int) (*models.Book, error) {
	var book models.Book
	err := r.DB.QueryRow(`SELECT id, title, author, year FROM books WHERE id = $1`, id).Scan(&book.ID, &book.Title, &book.Author, &book.BookYear)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) UpdateBook(book *models.Book) error {
	_, err := r.DB.Exec(`UPDATE books SET title = $1, author = $2, year = $3 WHERE id = $4`, book.Title, book.Author, book.BookYear, book.ID)
	return err
}

func (r *BookRepository) DeleteBook(id int) error {
	_, err := r.DB.Exec(`DELETE FROM books WHERE id = $1`, id)
	return err
}
