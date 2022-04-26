package main

import (
	"lan-chat/httpErrors"
	"lan-chat/logger"
	"net/http"
	"strings"
)

// StaticPageHandler to get static pages
func StaticPageHandler(w http.ResponseWriter, r *http.Request) { // To serve static pages
	fileLocation := strings.TrimPrefix(r.URL.Path, "/static/")
	logger.InfoLog.Println("static", fileLocation)
	http.ServeFile(w, r, "./templates/"+fileLocation)
}

func TemplateHandler(tempDir string) http.Handler { // tempDir : main template directory for a particular path
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			httpErrors.MethodNotAllowed(w)
			return
		}
		w.Header().Set("content-type", "text/html")
		fn := r.URL.Path
		logger.InfoLog.Println("template", fn)

		switch fn {
		case "/":
			http.ServeFile(w, r, tempDir+"/index.html")
		case "/favicon.ico":
			http.Redirect(w, r, "/static"+fn, http.StatusPermanentRedirect)
		default:
			file := fn + ".html"
			http.ServeFile(w, r, tempDir+file)
		}
	}
	return http.HandlerFunc(fn)
}
