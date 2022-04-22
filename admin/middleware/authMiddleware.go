package middleware

import (
	"context"
	"lan-chat/admin/jwt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
			username, err := jwt.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}
			ctxWithUser := context.WithValue(r.Context(), "username", username)
			r = r.WithContext(ctxWithUser)
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
