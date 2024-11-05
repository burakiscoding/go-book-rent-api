package api

import (
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/burakiscoding/go-book-rent/helpers"
	"github.com/burakiscoding/go-book-rent/store"
	"github.com/burakiscoding/go-book-rent/types"
	"github.com/gorilla/mux"
)

type BookHandler struct {
	store store.BookStore
}

func NewBookHandler(store store.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

func (h *BookHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) error {
	books, err := h.store.GetAll()
	if err != nil {
		return err
	}

	return helpers.WriteJSON(w, http.StatusOK, books)
}

func (h *BookHandler) HandleGetById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return helpers.InvalidRouteVariables()
	}

	book, err := h.store.GetById(id)
	if err != nil {
		return helpers.NotFoundData()
	}

	return helpers.WriteJSON(w, http.StatusOK, book)
}

func (h *BookHandler) HandleInsert(w http.ResponseWriter, r *http.Request) error {
	var book types.AddBookRequest
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		return helpers.InvalidJSON()
	}

	if book.Name == "" {
		return helpers.InvalidRequestData()
	}

	if err := h.store.Insert(book.Name); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}

func (h *BookHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return helpers.InvalidRouteVariables()
	}

	var book types.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		return helpers.InvalidJSON()
	}

	if book.Name == "" {
		return helpers.InvalidRequestData()
	}

	if err := h.store.Update(id, book.Name); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}

func (h *BookHandler) HandleDelete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return helpers.InvalidRouteVariables()
	}

	if err := h.store.Delete(id); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}
