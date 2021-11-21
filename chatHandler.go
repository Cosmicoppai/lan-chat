package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var userList [][]string

type User struct {
	Name string
}

// Message to receive msg from client
type Message struct {
	Message string
}

// SendMessage to broadcast message
type SendMessage struct {
	Name    string
	Message string
}

func (u *User) userName() string {
	return "sexy" + u.Name
}

type ChatServer interface {
	Listen(address string) error
	Broadcast(command interface{}) error
	Start()
	Close()
}

func chatHandler(conn net.Conn) {
	request := make([]byte, 1024)
	defer conn.Close()
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			log.Println(err)
			break
		}
		if readLen == 0 {
			// log.Println(conn.RemoteAddr().String(),"has left the chat")
			break // connection already closed by client
		} else if string(request[:readLen]) == "timestamp" {
			daytime := strconv.FormatInt(time.Now().Unix(), 10)
			_, _ = conn.Write([]byte(daytime))
		} else {
			_, _ = conn.Write([]byte(time.Now().String()))
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

func createSock(ip string) {
	service := ip + ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go chatHandler(conn)
	}
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
