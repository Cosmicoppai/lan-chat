package chat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type Message struct {
	typ string
	msg string
}

func Serve(conn net.Listener) error {
	for {
		req, e := conn.Accept()
		if e != nil {
			return e
		}
		go startChat(req)
	}
}

func startChat(conn net.Conn) {
	defer conn.Close()

	c := bufio.NewReader(conn)
	buff := make([]byte, 1024)      // make buffer of 1024 bytes
	r, _ := http.ReadRequest(c)     // to read header
	responseHeader := http.Header{} // initialize the response header
	rh := http.Response{}           // to write response
	if strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket" && r.Method == "GET" {

		key := r.Header.Get("Sec-WebSocket-Key")
		key = computeAcceptKey(key)
		responseHeader["Sec-WebSocket-Accept"] = []string{key}
		responseHeader["Upgrade"] = []string{"websocket"}
		responseHeader["Connection"] = []string{"Upgrade"}
		rh.StatusCode = http.StatusSwitchingProtocols
		rh.Header = responseHeader
		_ = rh.Write(conn) // write the response

		for {
			size, err := conn.Read(buff)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(string(buff[:size]))
			_, _ = conn.Write(buff[:size])
		}
	}
	// t := http.ResponseWriter(conn)
	_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
}
