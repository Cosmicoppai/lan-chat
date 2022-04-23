package middleware

import (
	"context"
	"lan-chat/admin/jwt"
	"lan-chat/httpErrors"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			token = strings.TrimPrefix(token, "Bearer ")
			claims, err := jwt.ValidateToken(token)
			if err != nil {
				if err == jwt.InvalidToken {
					httpErrors.UnAuthorized(w, "Invalid Token")
					return
				}
				httpErrors.InternalServerError(w)
				return
			}
			ctxWithUser := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctxWithUser)
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
