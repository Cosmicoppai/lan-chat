package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var EMAIL = "kanakchaudhari12@gmail.com" // email address to receive the suggestions

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_ = r.ParseForm() // parse the form
		movieName := r.Form.Get("movie_name")
		date := r.Form.Get("date")
		msg := r.Form.Get("msg")
		if movieName != "" && date != "" {
			message := fmt.Sprintf("%s has been requested on %s. <br> msg:- %s", movieName, date, msg)
			sendEmail(EMAIL, message)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Form Received! Thank you for your response"))
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return
}

func sendEmail(email string, msg string) int {
	dt := time.Now().Format("2006-01-02")
	requestUrl := "http://miraimail.herokuapp.com/"
	data := url.Values{}
	data.Add("Email", email)
	data.Add("Msg", msg)
	data.Add("Scheduled_date", dt)

	form, err := http.PostForm(requestUrl, data)
	if err != nil || form.StatusCode != 200 {
		body, _ := ioutil.ReadAll(form.Body)
		log.Println("Response", string(body))
		log.Println("Error", err)
		return -1
	}
	return 0
}
