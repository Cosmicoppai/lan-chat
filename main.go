package main

import (
	"github.com/cosmicoppai/lan-chat/chat"
	"log"
	"net"
	"net/http"
	"os/exec"
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
	log.Printf("server is listening on %s", ip)
	log.Fatalln(http.ListenAndServe(ip, server))
}

func getIpAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.To4().String()

}

func main() {
	cmd := exec.Command("cmd", "/c", `netsh wlan connect name="Laxmi 4"`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	log.Print(string(output))

	_ipAddress := getIpAddress()
	if _ipAddress != "" {
		go Server(_ipAddress) // start a goroutine
	} else {
		log.Fatalln("Pass a Valid IP Address")
	}

	conn, err_ := net.Listen("tcp", _ipAddress+":9000")
	if err_ != nil {
		log.Fatalln(err_)
	}
	e := chat.Serve(conn)
	if e != nil {
		return
	}

}
