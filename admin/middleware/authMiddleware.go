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

func AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuthorization(w, r)
		if err == nil {
			ctxWithUser := context.WithValue(r.Context(), "claims", claims)
			r = r.WithContext(ctxWithUser)
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func AdminMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		claims, err := checkAuthorization(w, r)
		if err == nil {
			if claims.IsAdmin {
				ctxWithUser := context.WithValue(r.Context(), "claims", claims)
				r = r.WithContext(ctxWithUser)
				next.ServeHTTP(w, r)
			} else {
				httpErrors.Forbidden(w)
				return
			}
		}
	}
	return http.HandlerFunc(fn)
}
