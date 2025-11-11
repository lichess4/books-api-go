// internal/transport/handlers.go
package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lichess4/books-api-go/internal/model"
)

// -------------- NUEVO: interfaz local para testear --------------
type BookService interface {
	GetAllBooks() ([]model.Book, error)
	CreateBook(model.Book) (model.Book, error)
	GetBookByID(id int) (model.Book, error)
	UpdateBook(id int, b model.Book) (model.Book, error)
	DeleteBook(id int) error
}

// ---------------------------------------------------------------

type BookHandler struct {
	service BookService // <- antes: *service.Service
}

// Cambia el constructor para aceptar la interfaz.
// Desde main puedes seguir pasando la implementación real *service.Service (porque implementa los métodos).
func New(s BookService) *BookHandler {
	return &BookHandler{service: s}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		books, err := h.service.GetAllBooks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(books)

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(created)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) HandleBookByID(w http.ResponseWriter, r *http.Request) {
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
		_ = json.NewEncoder(w).Encode(book)

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
		_ = json.NewEncoder(w).Encode(updated)

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
