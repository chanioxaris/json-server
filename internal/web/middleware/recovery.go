// Package middleware provides middleware support for REST API server.
package middleware

import (
	"net/http"
)

// Recovery is operating as middleware to handle any panic that may occur.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
