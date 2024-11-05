package api

import (
	"encoding/json"
	"net/http"

	"github.com/burakiscoding/go-book-rent/helpers"
	"github.com/burakiscoding/go-book-rent/store"
	"github.com/burakiscoding/go-book-rent/types"
)

type UserHandler struct {
	store store.UserStore
}

func NewUserHandler(store store.UserStore) *UserHandler {
	return &UserHandler{store: store}
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) error {
	var user types.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return helpers.InvalidJSON()
	}

	if user.Username == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
		return helpers.InvalidRequestData()
	}

	available, err := h.store.IsUsernameAvailable(user.Username)
	if err != nil {
		return err
	}

	if !available {
		return helpers.NewAPIError(http.StatusBadRequest, "username is already taken")
	}

	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}

	if err := h.store.Insert(user.Username, hashed, user.FirstName, user.LastName, types.RoleUser); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var user types.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return helpers.InvalidJSON()
	}

	if user.Username == "" || user.Password == "" {
		return helpers.InvalidRequestData()
	}

	foundUser, err := h.store.GetByUsername(user.Username)
	if err != nil {
		return helpers.BadCredentials()
	}

	isPasswordCorrect := helpers.CheckHashedPassword(foundUser.Password, user.Password)
	if !isPasswordCorrect {
		return helpers.BadCredentials()
	}

	token, err := helpers.CreateJWT(foundUser.Id, foundUser.Role)
	if err != nil {
		return err
	}

	return helpers.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) HandleGetDetails(w http.ResponseWriter, r *http.Request) error {
	tokenPayload, err := helpers.GetTokenPayloadFromContext(r)
	if err != nil {
		return err
	}

	user, err := h.store.GetById(tokenPayload.Id)
	if err != nil {
		return err
	}

	return helpers.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) HandleAdminRegister(w http.ResponseWriter, r *http.Request) error {
	var user types.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return helpers.InvalidJSON()
	}

	if user.Username == "" || user.Password == "" || user.FirstName == "" || user.LastName == "" {
		return helpers.InvalidRequestData()
	}

	available, err := h.store.IsUsernameAvailable(user.Username)
	if err != nil {
		return err
	}

	if !available {
		return helpers.NewAPIError(http.StatusBadRequest, "username is already taken")
	}

	hashed, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}

	if err := h.store.Insert(user.Username, hashed, user.FirstName, user.LastName, types.RoleAdmin); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}
