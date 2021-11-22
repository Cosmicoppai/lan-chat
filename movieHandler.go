package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, "./templates/main.html")
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	movie := strings.TrimPrefix(r.URL.Path, "/movie/")
	w.Header().Set("content-type", "video/mp4")
	w.Header().Set("accept-ranges", "bytes")
	http.ServeFile(w, r, "./movie/"+movie)
}

func currentMovies(w http.ResponseWriter, r *http.Request) {
	textData, _err := ioutil.ReadFile("./movie/schedule.txt") // get movie name from the file
	if _err != nil {
		http.Error(w, "some error occured", http.StatusInternalServerError)
		log.Println(_err)
		return
	}
	movie := string(textData)
	resp := make(map[string]string)
	resp["movie-name"] = movie
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Some Error Occured"))
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	return
}
