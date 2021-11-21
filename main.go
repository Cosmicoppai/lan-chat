package main

import (
	"log"
	"net"
	"net/http"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	// fileServer := http.FileServer(http.Dir("./templates"))
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, "./templates/main.html")
}

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

func Server(ip string) {
	server := http.NewServeMux()
	server.HandleFunc("/", StaticHandler)
	server.HandleFunc("/chat", chatHandler)     // endpoint for chat UI
	server.HandleFunc("/add-user", userHandler) // endpoint to add user
	log.Printf("server is listening on %s", ip)
	log.Fatalln(http.ListenAndServe(ip, server))
}

func getIpAddress() string {
	var ip string
	_addr, _err := net.InterfaceAddrs()
	if _err != nil {
		log.Println(_err)
		return ""
	}
	for _, a := range _addr {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip + ":80"

}

func main() {

	_ipAddress := getIpAddress()
	if _ipAddress != "" {
		Server(_ipAddress)
	} else {
		log.Fatalln("Pass a Valid IP Address")
	}

}
