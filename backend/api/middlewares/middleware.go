package middlewares

import (
	"net/http"

	"github.com/MinhL4m/blogs/api"
	"github.com/MinhL4m/blogs/api/auth"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)

		if err != nil {
			api.RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		next(w, r)
	}
}
