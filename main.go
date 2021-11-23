package main

import (
	"log"
	"net"
	"net/http"
)

func Server(ip string) {
	ip = ip + ":80"
	server := http.NewServeMux()
	server.HandleFunc("/", Home)
	server.HandleFunc("/static/", StaticPages)         // endpoint to get static pages
	server.HandleFunc("/send-suggestion", FormHandler) // to accept form-data
	server.HandleFunc("/movie_name", currentMovies)    // get the json response of current streaming movies
	server.HandleFunc("/get_movie/", GetMovie)         // endpoint to get movie
	server.HandleFunc("/get_sub/", GetSub)             // endpoint to get sub
	server.HandleFunc("/get_poster/", GetPoster)       //endpoint to  get current premiering movie poster
	// server.HandleFunc("/chat", chatHandler)     // endpoint for chat UI
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
	return ip

}

func main() {

	_ipAddress := getIpAddress()
	if _ipAddress != "" {
		Server(_ipAddress)
	} else {
		log.Fatalln("Pass a Valid IP Address")
	}

}
