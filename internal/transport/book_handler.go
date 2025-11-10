package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lichess4/books-api-go/internal/model"
	"github.com/lichess4/books-api-go/internal/service"
)

type BookHandler struct {
	service *service.Service
}

func New(s *service.Service) *BookHandler {
	return &BookHandler{
		service: s,
	}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books, err := h.service.GetAllBooks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Respond with the list of books (omitted for brevity)
		w.Header().Set("Content-Type", "application/json")
		// Assume we have a function toJSON to convert books to JSON
		json.NewEncoder(w).Encode(books)

	case http.MethodPost:
		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		created, err := h.service.CreateBook(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		// Respond with the created book (omitted for brevity)
		w.Header().Set("Content-Type", "application/json")
		// Assume we have a function toJSON to convert books to JSON
		json.NewEncoder(w).Encode(created)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL (omitted for brevity)
	// Assume we have a function extractID to get the ID from the URL
	idString := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		book, err := h.service.GetBookByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(book)

	case http.MethodPut:
		var book model.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		updated, err := h.service.UpdateBook(id, book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		if err := h.service.DeleteBook(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}
