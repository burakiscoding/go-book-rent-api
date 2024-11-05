package store

import (
	"database/sql"
	"time"

	"github.com/burakiscoding/go-book-rent/types"
	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Insert(username, password, firstName, lastName, role string) error {
	id := uuid.New()

	query := "INSERT INTO users (id, username, password, first_name, last_name, role, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, id, username, password, firstName, lastName, role, time.Now())

	return err
}

func (s *UserStore) GetByUsername(username string) (types.User, error) {
	var user types.User
	query := "SELECT id, username, password, first_name, last_name, role, created_at FROM users WHERE username = ?"
	err := s.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt)
	return user, err
}

func (s *UserStore) GetById(id string) (types.User, error) {
	var user types.User
	query := "SELECT id, username, password, first_name, last_name, role, created_at FROM users WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Password, &user.FirstName, &user.LastName, &user.Role, &user.CreatedAt)
	return user, err
}

func (s *UserStore) IsUsernameAvailable(username string) (bool, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)

	// Username is available
	if err == sql.ErrNoRows {
		return true, nil
	}

	return false, err
}
