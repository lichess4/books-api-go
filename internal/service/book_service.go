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

func (s *Service) GetAllBooks() ([]*model.Book, error) {

	//s.logger.Log("Fetching all books","")

	books, err := s.store.GetAll()
	if err != nil {
		//s.logger.Log("Error fetching books: %v\n", err.Error())
		return nil, err
	}

	return books, nil
}

func (s *Service) GetBookByID(id int) (*model.Book, error) {
	return s.store.GetByID(id)
}

func (s *Service) CreateBook(book model.Book) (*model.Book, error) {
	if book.Title == "" || book.Author == "" {
		return nil, errors.New("missing required fields")
	}

	return s.store.Create(&book)
}

func (s *Service) UpdateBook(id int, book model.Book) (*model.Book, error) {
	if book.Title == "" || book.Author == "" {
		return nil, errors.New("missing required fields")
	}
	return s.store.Update(id, &book)
}

func (s *Service) DeleteBook(id int) error {
	return s.store.Delete(id)
}
