package utils

import (
	"net/http"
	"regexp"
)

type Route struct {
	Method    string
	RegexPath *regexp.Regexp
	Handle    http.Handler
}

type CtxKey struct{}

func NewRoute(method string, pattern string, handler http.Handler) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(CtxKey{}).([]string)
	return fields[index]
}
