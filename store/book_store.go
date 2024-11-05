package store

import (
	"database/sql"

	"time"

	"github.com/burakiscoding/go-book-rent/types"
)

type BookStore struct {
	db *sql.DB
}

func NewBookStore(db *sql.DB) *BookStore {
	return &BookStore{
		db: db,
	}
}

func (s *BookStore) GetAll() ([]types.Book, error) {
	rows, err := s.db.Query("SELECT id, name, created_at, quantity FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []types.Book
	for rows.Next() {
		var b types.Book
		if err := rows.Scan(&b.Id, &b.Name, &b.CreatedAt, &b.Quantity); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (s *BookStore) GetById(id int) (types.Book, error) {
	var book types.Book
	query := "SELECT id, name, created_at, quantity FROM books WHERE id = ?"
	if err := s.db.QueryRow(query, id).Scan(&book.Id, &book.Name, &book.CreatedAt, &book.Quantity); err != nil {
		return types.Book{}, err
	}

	return book, nil
}

func (s *BookStore) Insert(name string) error {
	query := "INSERT INTO books (name, created_at) VALUES (?, ?)"
	if _, err := s.db.Exec(query, name, time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *BookStore) Update(id int, name string) error {
	query := "UPDATE books SET name = ? WHERE id = ?"
	if _, err := s.db.Exec(query, name, id); err != nil {
		return err
	}

	return nil
}

func (s *BookStore) Delete(id int) error {
	query := "DELETE FROM books	WHERE id = ?"
	if _, err := s.db.Exec(query, id); err != nil {
		return err
	}

	return nil
}
