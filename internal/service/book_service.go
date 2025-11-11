package service

import (
	"errors"

	"github.com/lichess4/books-api-go/internal/model"
	"github.com/lichess4/books-api-go/store"
)

/*type Logger interface {
	Log(message, err string)
}*/

type Service struct {
	store store.Store
	//logger Logger
}

func New(s store.Store) *Service {
	return &Service{
		store: s,
		//logger: nil,
	}
}

func (s *Service) GetAllBooks() ([]model.Book, error) {

	//s.logger.Log("Fetching all books","")

	books, err := s.store.GetAll()
	if err != nil {
		//s.logger.Log("Error fetching books: %v\n", err.Error())
		return nil, err
	}

	// Convertir punteros a valores
	result := make([]model.Book, len(books))
	for i, b := range books {
		if b != nil {
			result[i] = *b
		}
	}
	return result, nil
}

func (s *Service) GetBookByID(id int) (model.Book, error) {
	book, err := s.store.GetByID(id)
	if err != nil {
		return model.Book{}, err
	}
	if book == nil {
		return model.Book{}, errors.New("book not found")
	}
	return *book, nil
}

func (s *Service) CreateBook(book model.Book) (model.Book, error) {
	if book.Title == "" || book.Author == "" {
		return model.Book{}, errors.New("missing required fields")
	}

	created, err := s.store.Create(&book)
	if err != nil {
		return model.Book{}, err
	}
	if created == nil {
		return model.Book{}, errors.New("failed to create book")
	}
	return *created, nil
}

func (s *Service) UpdateBook(id int, book model.Book) (model.Book, error) {
	if book.Title == "" || book.Author == "" {
		return model.Book{}, errors.New("missing required fields")
	}
	updated, err := s.store.Update(id, &book)
	if err != nil {
		return model.Book{}, err
	}
	if updated == nil {
		return model.Book{}, errors.New("failed to update book")
	}
	return *updated, nil
}

func (s *Service) DeleteBook(id int) error {
	return s.store.Delete(id)
}
