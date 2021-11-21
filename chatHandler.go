package main

import (
	"log"
	"net/http"
)

var userList [][]string

type User struct {
	name string
}

func (u *User) userName() string {
	return "sexy" + u.name
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := w.Write([]byte("よこそう チャットルーム へ"))
		if err != nil {
			log.Fatalln(err)
			return
		}
	}

}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_username := (&User{r.PostForm.Get("username")}).userName()
		userinfo := []string{_username, "ip"}
		userList = append(userList, userinfo)
		return

	}
	http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
	return
}
