package utils

import (
	"context"
	"lan-chat/httpErrors"
	"net/http"
	"strings"
)

func Router(route []Route) http.Handler {
	_routes := route
	fn := func(w http.ResponseWriter, r *http.Request) {
		var allow []string
		for _, route := range _routes {
			matches := route.RegexPath.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 { // if any path matches
				if r.Method != route.Method { // if request Method is different, continue the url match
					allow = append(allow, route.Method)
					continue
				}
				// if url with exact path parameter and Method is found, call the handler with updated request context
				ctx := context.WithValue(r.Context(), CtxKey{}, matches[1:])
				route.Handle.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		if len(allow) > 0 {
			w.Header().Set("Allow", strings.Join(allow, ", "))
			httpErrors.MethodNotAllowed(w)
			return
		}
		httpErrors.NotFound(w)
	}
	return http.HandlerFunc(fn)
}
