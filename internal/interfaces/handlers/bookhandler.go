package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dias221467/MicroServices/internal/domain/models"
	"github.com/Dias221467/MicroServices/internal/interfaces"
	"github.com/gorilla/mux"
)

type bookHandler struct {
	bookUsecase interfaces.BookUsecase
}

func NewBookHandler(bookUsecase interfaces.BookUsecase) interfaces.BookHandler {
	return &bookHandler{bookUsecase}
}

func (h *bookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" || book.BookYear == 0 {
		http.Error(w, "Title, author, and year are required", http.StatusBadRequest)
		return
	}

	if err := h.bookUsecase.AddBook(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *bookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookUsecase.GetBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *bookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := h.bookUsecase.GetBookByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *bookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if book.Title == "" || book.Author == "" || book.BookYear == 0 {
		http.Error(w, "Title, author, and year are required", http.StatusBadRequest)
		return
	}

	book.ID = id
	if err := h.bookUsecase.UpdateBook(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *bookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.bookUsecase.DeleteBook(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
