package helpers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/burakiscoding/go-book-rent/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewAPIError(status int, message string) APIError {
	return APIError{Status: status, Message: message}
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %s and status: %d", e.Message, e.Status)
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, "invalid JSON request data")
}

func InvalidRequestData() APIError {
	return NewAPIError(http.StatusBadRequest, "invalid request data")
}

func NotFoundData() APIError {
	return NewAPIError(http.StatusNotFound, "data not found")
}

func InvalidRouteVariables() APIError {
	return NewAPIError(http.StatusBadRequest, "invalid route variables")
}

func BadCredentials() APIError {
	return NewAPIError(http.StatusUnauthorized, "bad credentials")
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteOK(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{"message": "Mission Completed"})
}

func MakeHandler(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				WriteJSON(w, apiErr.Status, apiErr)
			} else {
				errServer := NewAPIError(http.StatusInternalServerError, "internal server error")
				WriteJSON(w, errServer.Status, errServer)
			}
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckHashedPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func CreateJWT(id, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"role": role,
		"iss":  "book-rent",
		"aud":  "normal",
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"iat":  time.Now().Unix(),
	})

	return claims.SignedString(secret)
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	substrings := strings.Split(authHeader, " ")

	if len(substrings) != 2 {
		return "", BadCredentials()
	}
	if substrings[0] != "Bearer" {
		return "", BadCredentials()
	}

	token := substrings[1]
	return token, nil
}

func GetTokenPayload(tokenString string) (types.TokenPayload, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return types.TokenPayload{}, BadCredentials()
	}

	if !token.Valid {
		return types.TokenPayload{}, BadCredentials()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return types.TokenPayload{}, BadCredentials()
	}

	return types.TokenPayload{Id: claims["sub"].(string), Role: claims["role"].(string)}, nil
}

func GetTokenPayloadFromContext(r *http.Request) (types.TokenPayload, error) {
	id := r.Context().Value(types.KeyId).(string)
	role := r.Context().Value(types.KeyRole).(string)

	if id == "" || role == "" {
		return types.TokenPayload{}, BadCredentials()
	}

	return types.TokenPayload{Id: id, Role: role}, nil
}
