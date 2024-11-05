package types

import "time"

type ContextKey string

const (
	RoleUser          string     = "user"
	RoleAdmin         string     = "admin"
	KeyId             ContextKey = "KeyId"
	KeyRole           ContextKey = "KeyRole"
	MinRentTimeInDays int        = 1
	MaxRentTimeInDays int        = 30
)

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Book struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Quantity  int       `json:"quantity"`
}

type RentHistory struct {
	Id                 string     `json:"id"`
	BookId             int        `json:"book_id"`
	UserId             string     `json:"user_id"`
	RentStartTime      time.Time  `json:"rent_start_time"`
	RentReturnTime     *time.Time `json:"rent_return_time"`
	RentDurationInDays int        `json:"rent_duration_in_days"`
}

type AddBookRequest struct {
	Name string `json:"name"`
}

type UpdateBookRequest struct {
	Name string `json:"name"`
}

type RegisterUserRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenPayload struct {
	Id   string
	Role string
}

type RentBookRequest struct {
	BookId         int `json:"book_id"`
	DurationInDays int `json:"duration_in_days"`
}

type ReturnBookRequest struct {
	Id string `json:"id"`
}

type UserRentHistory struct {
	Id                 string     `json:"id"`
	RentStartTime      time.Time  `json:"rent_start_time"`
	RentReturnTime     *time.Time `json:"rent_return_time"`
	RentDurationInDays int        `json:"rent_duration_in_days"`
	BookName           string     `json:"book_name"`
}
