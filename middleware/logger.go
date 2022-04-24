package middleware

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s- %s  %s\n", r.RemoteAddr, r.RequestURI, r.Method)
			next.ServeHTTP(w, r)
		})
}
