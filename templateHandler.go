package main

import (
	"lan-chat/httpErrors"
	"net/http"
	"strings"
)

func StaticPageHandler(w http.ResponseWriter, r *http.Request) { // To serve static pages
	fileLocation := strings.TrimPrefix(r.URL.Path, "/static/")
	http.ServeFile(w, r, "./templates/"+fileLocation)
}

func TemplateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpErrors.MethodNotAllowed(w)
		return
	}
	w.Header().Set("content-type", "text/html")
	fn := r.URL.Path

	switch fn {
	case "/":
		http.ServeFile(w, r, "./templates/index.html")
	case "/favicon.ico":
		http.Redirect(w, r, "/static"+fn, 301)
	default:
		file := fn + ".html"
		http.ServeFile(w, r, "./templates/"+file)
	}
}
