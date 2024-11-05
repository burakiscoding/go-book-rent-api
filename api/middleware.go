package api

import (
	"context"
	"net/http"

	"github.com/burakiscoding/go-book-rent/helpers"
	"github.com/burakiscoding/go-book-rent/types"
)

func HandleAuth(f helpers.APIFunc) helpers.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenString, err := helpers.GetTokenFromHeader(r)
		if err != nil {
			return err
		}

		tokenPayload, err := helpers.GetTokenPayload(tokenString)
		if err != nil {
			return err
		}

		if tokenPayload.Role != types.RoleUser && tokenPayload.Role != types.RoleAdmin {
			return helpers.BadCredentials()
		}

		ctx := context.WithValue(r.Context(), types.KeyId, tokenPayload.Id)
		ctx = context.WithValue(ctx, types.KeyRole, tokenPayload.Role)
		return f(w, r.WithContext(ctx))
	}
}

func HandleAdminAuth(f helpers.APIFunc) helpers.APIFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenString, err := helpers.GetTokenFromHeader(r)
		if err != nil {
			return err
		}

		tokenPayload, err := helpers.GetTokenPayload(tokenString)
		if err != nil {
			return err
		}

		if tokenPayload.Role != types.RoleAdmin {
			return helpers.BadCredentials()
		}

		ctx := context.WithValue(r.Context(), types.KeyId, tokenPayload.Id)
		ctx = context.WithValue(ctx, types.KeyRole, tokenPayload.Role)
		return f(w, r.WithContext(ctx))
	}
}
