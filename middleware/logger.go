package middleware

import (
	"lan-chat/logger"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.InfoLog.Printf("%s <- %s\n", r.RequestURI, r.Method)
			next.ServeHTTP(w, r)
		})
}
