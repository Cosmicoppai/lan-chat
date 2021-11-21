package main

import (
	"log"
	"net"
	"net/http"
)

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
