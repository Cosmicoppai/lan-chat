package main

import (
	"flag"
	"lan-chat/admin"
	"lan-chat/admin/middleware"
	"lan-chat/admin/users"
	"lan-chat/chat"
	"lan-chat/movieHandler"
	"lan-chat/suggestions"
	"log"
	"net"
	"net/http"
	"os"
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
	server.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(users.ListUsers)))
	server.Handle("/user", middleware.AuthMiddleware(http.HandlerFunc(users.Handler)))
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
	if len(os.Args) < 2 {
		log.Fatalln("Command not specified !..\nExpected runserver or create-superuser sub-command")
	}
	_ipAddress := getIpAddress()
	if _ipAddress == "" {
		os.Exit(1)
	}

	envFilePath, _ := filepath.Abs(".env")
	admin.LoadEnv(envFilePath)
	admin.InitializeDB()
	defer admin.Db.Close()

	runServerCmd := flag.NewFlagSet("runserver", flag.ExitOnError)
	host := runServerCmd.String("host", _ipAddress, "Host name or IP address")
	socketPort := runServerCmd.String("sockPort", "9000", "Port address of socket server")

	superUserCmd := flag.NewFlagSet("create-superuser", flag.ExitOnError)

	username := superUserCmd.String("username", "", "Username")
	password := superUserCmd.String("password", "", "Password")
	switch os.Args[1] {
	case "create-superuser":
		handleCreateSuperUser(superUserCmd, username, password)
	case "runserver":
		handleRunServer(runServerCmd, host, socketPort)
	default:
		log.Fatalf("unknown command %s\n", os.Args[1])

	}

}

func handleRunServer(runServerCmd *flag.FlagSet, host *string, sockPort *string) {
	err := runServerCmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatalln(err)
	}
	runServers(*host, *sockPort)

}

func handleCreateSuperUser(superUserCmd *flag.FlagSet, username *string, password *string) {
	err := superUserCmd.Parse(os.Args[2:])
	if err != nil || *username == "" || *password == "" {
		if err == nil {
			log.Fatalln("Enter both username and password")
		}
		log.Fatalln(err)
	}
	users.CreateSuperUser(*username, *password)

}

func runServers(host string, socketPort string) {

	go Server(host) // start a goroutine

	conn, err_ := net.Listen("tcp", host+":"+socketPort)
	if err_ != nil {
		log.Fatalln(err_)
	}
	log.Printf("Socket server is listening on %s:%s", host, socketPort)
	e := chat.Serve(conn)
	if e != nil {
		log.Println("error in chat.Serve", e)
		return
	}
}
