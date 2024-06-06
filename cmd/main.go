package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/Dias221467/MicroServices/internal/domain/models"
	adapters "github.com/Dias221467/MicroServices/internal/interfaces/adapters/postgres"
	"github.com/Dias221467/MicroServices/internal/usecases"
	pb "github.com/Dias221467/MicroServices/proto" // Import the generated protobuf code
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	// Initialize the logger

	db, err := sql.Open("postgres", "postgresql://postgres:lbfc2005@localhost:5432/microservices?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	bookRepo := adapters.NewBookRepository(db)
	bookUsecase := usecases.NewBookUsecase(bookRepo)

	r := mux.NewRouter()
	r.HandleFunc("/books", createBookHandler(bookUsecase)).Methods("POST")
	r.HandleFunc("/books", getBooksHandler(bookUsecase)).Methods("GET")
	r.HandleFunc("/books/{id}", getBookHandler(bookUsecase)).Methods("GET")
	r.HandleFunc("/books/{id}", updateBookHandler(bookUsecase)).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBookHandler(bookUsecase)).Methods("DELETE")

	// Start the server
	logger.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createBookHandler(usecase *usecases.BookUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := usecase.AddBook(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func getBooksHandler(usecase *usecases.BookUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := usecase.GetBooks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(books)
	}
}

func getBookHandler(usecase *usecases.BookUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		book, err := usecase.GetBookByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

func updateBookHandler(usecase *usecases.BookUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		book.ID = id
		if err := usecase.UpdateBook(&book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

func deleteBookHandler(usecase *usecases.BookUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := usecase.DeleteBook(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type server struct {
	usecase *usecases.BookUsecase
	pb.UnimplementedBookServiceServer
}

func NewBookServiceServer(usecase *usecases.BookUsecase) pb.BookServiceServer {
	return &server{usecase: usecase}
}

func runGRPCServer(usecase *usecases.BookUsecase) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatalf("Failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookServiceServer(s, NewBookServiceServer(usecase))

	logger.Println("Starting gRPC server on port 50051...")
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
