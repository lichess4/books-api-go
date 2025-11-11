package transport_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lichess4/books-api-go/internal/model"
	"github.com/lichess4/books-api-go/internal/transport"
)

// ---------- Fake Service para tests ----------
type fakeService struct {
	data map[int]model.Book
	next int
	// flags de error opcionales
	errList   error
	errCreate error
	errGet    error
	errUpdate error
	errDelete error
}

func newFake() *fakeService {
	return &fakeService{
		data: map[int]model.Book{
			1: {ID: 1, Title: "Clean Code", Author: "Robert C. Martin"},
			2: {ID: 2, Title: "The Pragmatic Programmer", Author: "Andrew Hunt"},
		},
		next: 3,
	}
}

func (f *fakeService) GetAllBooks() ([]model.Book, error) {
	if f.errList != nil {
		return nil, f.errList
	}
	out := make([]model.Book, 0, len(f.data))
	for _, b := range f.data {
		out = append(out, b)
	}
	return out, nil
}

func (f *fakeService) CreateBook(b model.Book) (model.Book, error) {
	if f.errCreate != nil {
		return model.Book{}, f.errCreate
	}
	b.ID = f.next
	f.next++
	f.data[b.ID] = b
	return b, nil
}

func (f *fakeService) GetBookByID(id int) (model.Book, error) {
	if f.errGet != nil {
		return model.Book{}, f.errGet
	}
	if b, ok := f.data[id]; ok {
		return b, nil
	}
	return model.Book{}, errors.New("not found")
}

func (f *fakeService) UpdateBook(id int, b model.Book) (model.Book, error) {
	if f.errUpdate != nil {
		return model.Book{}, f.errUpdate
	}
	if _, ok := f.data[id]; !ok {
		return model.Book{}, errors.New("not found")
	}
	b.ID = id
	f.data[id] = b
	return b, nil
}

func (f *fakeService) DeleteBook(id int) error {
	if f.errDelete != nil {
		return f.errDelete
	}
	if _, ok := f.data[id]; !ok {
		return errors.New("not found")
	}
	delete(f.data, id)
	return nil
}

// ---------- Helpers ----------
func mustJSON(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewReader(b)
}

// ================== TESTS =====================

func TestHandleBooks_GET_OK(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBooks).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var got []model.Book
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d want=2", len(got))
	}
}

func TestHandleBooks_POST_Created(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	body := mustJSON(t, map[string]any{
		"title":  "DDD",
		"author": "Eric Evans",
	})
	req := httptest.NewRequest(http.MethodPost, "/books", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	http.HandlerFunc(h.HandleBooks).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d want=%d body=%s", w.Code, http.StatusCreated, w.Body.String())
	}

	var got model.Book
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.ID == 0 || got.Title != "DDD" || got.Author != "Eric Evans" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestHandleBooks_MethodNotAllowed(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodPatch, "/books", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBooks).ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d want=%d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandleBookByID_GET_Found(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodGet, "/books/1", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d", w.Code, http.StatusOK)
	}

	var got model.Book
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.ID != 1 {
		t.Fatalf("id=%d want=1", got.ID)
	}
}

func TestHandleBookByID_GET_NotFound(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d want=%d body=%s", w.Code, http.StatusNotFound, w.Body.String())
	}
}

func TestHandleBookByID_GET_BadID(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want=%d", w.Code, http.StatusBadRequest)
	}
}

func TestHandleBookByID_PUT_Update(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	body := mustJSON(t, map[string]any{
		"title":  "Clean Code (2nd Ed.)",
		"author": "Robert C. Martin",
	})
	req := httptest.NewRequest(http.MethodPut, "/books/1", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=%d body=%s", w.Code, http.StatusOK, w.Body.String())
	}

	var got model.Book
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.ID != 1 || got.Title != "Clean Code (2nd Ed.)" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestHandleBookByID_DELETE_NoContent(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodDelete, "/books/2", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status=%d want=%d", w.Code, http.StatusNoContent)
	}

	// Asegura que ya no exista
	req2 := httptest.NewRequest(http.MethodGet, "/books/2", nil)
	w2 := httptest.NewRecorder()
	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w2, req2)
	if w2.Code != http.StatusNotFound {
		t.Fatalf("status after delete=%d want=%d", w2.Code, http.StatusNotFound)
	}
}

func TestHandleBookByID_MethodNotAllowed(t *testing.T) {
	svc := newFake()
	h := transport.New(svc)

	req := httptest.NewRequest(http.MethodPatch, "/books/1", nil)
	w := httptest.NewRecorder()

	http.HandlerFunc(h.HandleBookByID).ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d want=%d", w.Code, http.StatusMethodNotAllowed)
	}
}
