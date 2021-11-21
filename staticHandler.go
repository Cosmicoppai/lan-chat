package main

import (
	"net/http"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, "./templates/main.html")
}
