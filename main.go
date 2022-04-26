package main

import (
	"flag"
	"lan-chat/admin"
	"lan-chat/admin/users"
	"lan-chat/chat"
	"lan-chat/logger"
	"lan-chat/middleware"
	"lan-chat/utils"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func Server(ip string) {
	ip = ip + ":80"
	if len(AppRoutes) > 0 {
		for _, routes := range AppRoutes {
			Routes = append(Routes, routes...)
		}
	}

	muxWithLogging := middleware.Logger(utils.Router(Routes))
	logger.InfoLog.Printf("Http server is listening on %s", ip)

	srv := &http.Server{
		ErrorLog: logger.ErrorLog,
		Handler:  muxWithLogging,
		Addr:     ip,
	}
	logger.ErrorLog.Fatalln(srv.ListenAndServe())
}

func getIpAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.ErrorLog.Fatalln("Make Sure you're connected to the internet")
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	return conn.LocalAddr().(*net.UDPAddr).IP.To4().String()

}

func main() {
	if len(os.Args) < 2 {
		logger.ErrorLog.Fatalln("Command not specified !..\nExpected runserver or create-superuser sub-command")
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
		logger.ErrorLog.Fatalf("unknown command %s\n", os.Args[1])

	}

}

func handleRunServer(runServerCmd *flag.FlagSet, host *string, sockPort *string) {
	err := runServerCmd.Parse(os.Args[2:])
	if err != nil {
		logger.ErrorLog.Fatalln(err)
	}
	runServers(*host, *sockPort)

}

func handleCreateSuperUser(superUserCmd *flag.FlagSet, username *string, password *string) {
	err := superUserCmd.Parse(os.Args[2:])
	if err != nil || *username == "" || *password == "" {
		if err == nil {
			logger.ErrorLog.Fatalln("Enter both username and password")
		}
		logger.ErrorLog.Fatalln(err)
	}
	users.CreateSuperUser(*username, *password)

}

func runServers(host string, socketPort string) {

	go Server(host) // start a goroutine

	conn, err_ := net.Listen("tcp", host+":"+socketPort)
	if err_ != nil {
		logger.ErrorLog.Fatalln(err_)
	}
	logger.InfoLog.Printf("Socket server is listening on %s:%s", host, socketPort)
	e := chat.Serve(conn)
	if e != nil {
		logger.ErrorLog.Println("error in chat.Serve", e)
		return
	}
}
