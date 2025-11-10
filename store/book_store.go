package store

import (
	"database/sql"

	"github.com/lichess4/books-api-go/internal/model"
)

type Store interface {
	GetAll() ([]*model.Book, error)
	GetByID(id int) (*model.Book, error)
	Create(book *model.Book) (*model.Book, error)
	Update(id int, book *model.Book) (*model.Book, error)
	Delete(id int) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetAll() ([]*model.Book, error) {
	q := `Select id, title, author from books`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {

		b := model.Book{}

		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, &b)
	}
	return books, nil
}

func (s *store) GetByID(id int) (*model.Book, error) {

	q := `Select id, title, author from books where id = ?`

	b := model.Book{}

	err := s.db.QueryRow(q, id).Scan(&b.ID, &b.Title, &b.Author)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *store) Create(book *model.Book) (*model.Book, error) {
	q := `Insert into books (title, author) values (?, ?)`
	resp, err := s.db.Exec(q, book.Title, book.Author)
	if err != nil {
		return nil, err
	}
	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}
	book.ID = int(id)
	return book, nil
}

func (s *store) Update(id int, book *model.Book) (*model.Book, error) {
	q := `Update books set title = ?, author = ? where id = ?`
	_, err := s.db.Exec(q, book.Title, book.Author, id)
	if err != nil {
		return nil, err
	}

	book.ID = id
	return book, nil
}

func (s *store) Delete(id int) error {
	q := `Delete from books where id = ?`
	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}
