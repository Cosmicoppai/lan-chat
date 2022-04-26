package middleware

import (
	"context"
	"lan-chat/admin/jwt"
	"lan-chat/httpErrors"
	"net/http"
	"strings"
)

func checkAuthorization(w http.ResponseWriter, r *http.Request) (jwt.Claims, error) {
	token := r.Header.Get("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			if err == jwt.InvalidToken {
				httpErrors.UnAuthorized(w, "Invalid Token")
				return jwt.Claims{}, err
			}
			httpErrors.InternalServerError(w)
			return jwt.Claims{}, err
		}
		return claims, nil

	}
	httpErrors.UnAuthorized(w, "Invalid Token")
	return jwt.Claims{}, jwt.InvalidToken
}

func AuthMiddleware(next func(http.ResponseWriter, *http.Request)) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuthorization(w, r)
		if err == nil {
			ctxWithUser := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctxWithUser)
			http.HandlerFunc(next).ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func AdminMiddleware(next func(http.ResponseWriter, *http.Request)) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuthorization(w, r)
		if err != nil && claims.IsAdmin {
			httpErrors.Forbidden(w)
			return
		}
		ctxWithUser := context.WithValue(r.Context(), "claims", claims)
		r = r.WithContext(ctxWithUser)
		http.HandlerFunc(next).ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
