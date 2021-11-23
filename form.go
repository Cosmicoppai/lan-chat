package main

import (
	"fmt"
	"net/http"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_ = r.ParseForm() // parse the form
		movieName := r.Form.Get("movie_name")
		date := r.Form.Get("date")
		msg := r.Form.Get("msg")
		if movieName != "" && date != "" {
			fmt.Printf("%s requested on %s and has sent message %s\n", movieName, date, msg)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Form Received! Thank you for your response"))
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return
}
