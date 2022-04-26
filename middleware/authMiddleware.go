package middleware

import (
	"context"
	"lan-chat/admin/jwt"
	"lan-chat/httpErrors"
	"net/http"
	"strings"
)

func checkAuthorization(r *http.Request) (jwt.Claims, error) {
	token := r.Header.Get("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := jwt.ValidateToken(token)
		return claims, err

	}
	return jwt.Claims{}, jwt.InvalidToken
}

func AuthMiddleware(next func(http.ResponseWriter, *http.Request)) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuthorization(r)
		if authError(w, err) {
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctxWithUser)
		http.HandlerFunc(next).ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func AdminMiddleware(next func(http.ResponseWriter, *http.Request)) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuthorization(r)
		if authError(w, err) {
			return
		}
		if !claims.IsAdmin {
			httpErrors.Forbidden(w)
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctxWithUser)
		http.HandlerFunc(next).ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func authError(w http.ResponseWriter, err error) bool {
	if err != nil {
		if err == jwt.InvalidToken {
			httpErrors.UnAuthorized(w, "Invalid Token")
			return true
		}
		httpErrors.InternalServerError(w)
		return true
	}
	return false
}
