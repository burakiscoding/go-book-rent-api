package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/burakiscoding/go-book-rent/types"
	"github.com/google/uuid"
)

type RentStore struct {
	db *sql.DB
}

func NewRentStore(db *sql.DB) *RentStore {
	return &RentStore{db: db}
}

func (s *RentStore) GetAllHistory() ([]types.RentHistory, error) {
	rows, err := s.db.Query("SELECT id, book_id, user_id, rent_duration_in_days, rent_start_time, rent_return_time FROM book_rent_history")
	if err != nil {
		return nil, err
	}

	var history []types.RentHistory
	for rows.Next() {
		var h types.RentHistory
		if err := rows.Scan(&h.Id, &h.BookId, &h.UserId, &h.RentDurationInDays, &h.RentStartTime, &h.RentReturnTime); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (s *RentStore) GetHistoryById(id string) (types.RentHistory, error) {
	var h types.RentHistory
	query := "SELECT id, book_id, user_id, rent_start_time, rent_return_time, rent_duration_in_days FROM book_rent_history WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(&h.Id, &h.BookId, &h.UserId, &h.RentStartTime, &h.RentReturnTime, &h.RentDurationInDays)
	if err != nil {
		return types.RentHistory{}, nil
	}

	return h, nil
}

func (s *RentStore) GetUserHistory(userId string) ([]types.UserRentHistory, error) {
	query := "SELECT R.id, R.rent_start_time, R.rent_return_time, R.rent_duration_in_days, B.name FROM book_rent_history AS R INNER JOIN books AS B on R.book_id = B.id WHERE R.user_id = ?"
	rows, err := s.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var history []types.UserRentHistory
	for rows.Next() {
		var h types.UserRentHistory
		if err := rows.Scan(&h.Id, &h.RentStartTime, &h.RentReturnTime, &h.RentDurationInDays, &h.BookName); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

func (s *RentStore) RentBook(ctx context.Context, bookId int, userId string, durationInDays int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert new record to the book_rent_history table
	id := uuid.New()
	query := "INSERT INTO book_rent_history (id, book_id, user_id, rent_duration_in_days, rent_start_time) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, query, id, bookId, userId, durationInDays, time.Now())
	if err != nil {
		return err
	}

	// Decrase the quantity in the books table
	query = "UPDATE books SET quantity = quantity - 1 WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, bookId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *RentStore) ReturnBook(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Find the book
	var bookId int
	query := "SELECT book_id FROM book_rent_history WHERE id = ?"
	err = tx.QueryRowContext(ctx, query, id).Scan(&bookId)
	if err != nil {
		return err
	}

	// Update rent_return_time in the rent_book_history table
	query = "UPDATE book_rent_history SET rent_return_time = ? WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	// Increase the quantity in the book table
	query = "UPDATE books SET quantity = quantity + 1 WHERE id = ?"
	_, err = tx.ExecContext(ctx, query, bookId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
