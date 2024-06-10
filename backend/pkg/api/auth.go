package api

import (
	"github.com/d-rk/checkin-system/pkg/auth"
	"golang.org/x/net/context"
	"net/http"
)

const AuthenticatedUserID = "auth.UserID"

func AuthMiddleware() MiddlewareFunc {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if bearerTokenRequired(r) {

				token, err := auth.FindToken(r)
				if err != nil {
					handlerError(w, r, InvalidTokenErr.Wrap(err))
					return
				}

				claims, err := auth.ValidateToken(token)
				if err != nil {
					handlerError(w, r, InvalidTokenErr.Wrap(err))
					return
				}

				r = r.WithContext(context.WithValue(r.Context(), AuthenticatedUserID, claims.UserID))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func bearerTokenRequired(r *http.Request) bool {
	_, ok := r.Context().Value(BearerAuthScopes).([]string)
	return ok
}
