package main

import (
	"lan-chat/admin"
	"lan-chat/admin/users"
	"lan-chat/chat"
	"lan-chat/movieHandler"
	"lan-chat/suggestions"
	"log"
	"net"
	"net/http"
	"path/filepath"
)

func Server(ip string) {
	ip = ip + ":80"
	server := http.NewServeMux()
	server.HandleFunc("/", movieHandler.Home)
	server.HandleFunc("/static/", movieHandler.StaticPages)        // endpoint to get static pages
	server.HandleFunc("/send-suggestion", suggestions.FormHandler) // to accept form-data
	server.HandleFunc("/list-movies", movieHandler.ListVideos)     // get the json response of current streaming movies
	server.HandleFunc("/file/", movieHandler.GetFile)              // endpoint to get movie
	server.HandleFunc("/user", users.Handler)
	server.HandleFunc("/login", users.Login)
	// server.HandleFunc("/bwahahaha/", admin.Handler)
	log.Printf("Http server is listening on %s", ip)
	log.Fatalln(http.ListenAndServe(ip, server))
}

func getIpAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println("Make Sure you're connected to the internet")
		log.Panic(err)
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	return conn.LocalAddr().(*net.UDPAddr).IP.To4().String()

}

func main() {
	envFilePath, _ := filepath.Abs(".env")
	admin.LoadEnv(envFilePath)
	admin.InitializeDB()
	defer admin.Db.Close()

	_ipAddress := getIpAddress()
	if _ipAddress != "" {
		go Server(_ipAddress) // start a goroutine
	} else {
		return
	}

	conn, err_ := net.Listen("tcp", _ipAddress+":9000")
	if err_ != nil {
		log.Fatalln(err_)
	}
	log.Printf("Socket server is listening on %s", _ipAddress+":9000")
	e := chat.Serve(conn)
	if e != nil {
		log.Println("error in chat.Serve", e)
		return
	}

}
