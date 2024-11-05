package api

import (
	"encoding/json"
	"net/http"

	"github.com/burakiscoding/go-book-rent/helpers"
	"github.com/burakiscoding/go-book-rent/store"
	"github.com/burakiscoding/go-book-rent/types"
)

type RentHandler struct {
	store     store.RentStore
	bookStore store.BookStore
}

func NewRentHandler(store store.RentStore, bookStore store.BookStore) *RentHandler {
	return &RentHandler{store: store, bookStore: bookStore}
}

func (h *RentHandler) HandleRentBook(w http.ResponseWriter, r *http.Request) error {
	tokenPayload, err := helpers.GetTokenPayloadFromContext(r)
	if err != nil {
		return err
	}

	var request types.RentBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.InvalidJSON()
	}

	if request.BookId == 0 ||
		request.DurationInDays < types.MinRentTimeInDays ||
		request.DurationInDays > types.MaxRentTimeInDays {
		return helpers.InvalidRequestData()
	}

	book, err := h.bookStore.GetById(request.BookId)
	if err != nil {
		return err
	}

	if book.Quantity <= 0 {
		return helpers.NotFoundData()
	}

	if err := h.store.RentBook(r.Context(), request.BookId, tokenPayload.Id, request.DurationInDays); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}

func (h *RentHandler) HandleReturnBook(w http.ResponseWriter, r *http.Request) error {
	tokenPayload, err := helpers.GetTokenPayloadFromContext(r)
	if err != nil {
		return err
	}

	var request types.ReturnBookRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return helpers.InvalidJSON()
	}

	if request.Id == "" {
		return helpers.InvalidRequestData()
	}

	history, err := h.store.GetHistoryById(request.Id)
	if err != nil {
		return err
	}

	// It's not your rent history
	if history.UserId != tokenPayload.Id {
		return helpers.BadCredentials()
	}

	// Book already returned
	if history.RentReturnTime != nil {
		return helpers.NotFoundData()
	}

	if err := h.store.ReturnBook(r.Context(), request.Id); err != nil {
		return err
	}

	return helpers.WriteOK(w)
}

func (h *RentHandler) HandleGetAllHistory(w http.ResponseWriter, r *http.Request) error {
	history, err := h.store.GetAllHistory()
	if err != nil {
		return err
	}

	return helpers.WriteJSON(w, http.StatusOK, history)
}

func (h *RentHandler) HandleGetUserHistory(w http.ResponseWriter, r *http.Request) error {
	tokenPayload, err := helpers.GetTokenPayloadFromContext(r)
	if err != nil {
		return err
	}

	history, err := h.store.GetUserHistory(tokenPayload.Id)
	if err != nil {
		return err
	}

	return helpers.WriteJSON(w, http.StatusOK, history)
}
