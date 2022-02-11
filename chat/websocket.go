package chat

import (
	"bufio"
	"net"
	"net/http"
	"strings"
)

const (
	finalBit = 1 << 7 // 10000000  set the first bit
	// rsv1Bit  = 1 << 6
	// rsv2Bit  = 1 << 5
	// rsv3Bit  = 1 << 4

	maskBit        = 1 << 7 // 10000000
	maskKeyLen int = 4      // length of maskKeyLen
	// maxFrameHeaderSize     = 2 + 8 + 4

	continuationFrame = 0
	// noFrame           = -1
)

const (
	TextMessage = 1
	PingMessage = 9
	PongMessage = 10
)

func connectionUpgrade(conn net.Conn) (success bool) {

	c := bufio.NewReader(conn)
	r, _ := http.ReadRequest(c)     // to read header
	responseHeader := http.Header{} // initialize the response header
	rh := http.Response{}           // to write response
	if strings.ToLower(r.Header.Get("Connection")) == "upgrade" && r.Header.Get("Sec-WebSocket-Key") != "" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket" && r.Method == "GET" {

		key := r.Header.Get("Sec-WebSocket-Key")
		key = computeAcceptKey(key)
		responseHeader["Sec-WebSocket-Accept"] = []string{key}
		responseHeader["Upgrade"] = []string{"websocket"}
		responseHeader["Connection"] = []string{"Upgrade"}
		rh.StatusCode = http.StatusSwitchingProtocols
		rh.Header = responseHeader // Append the Header
		_ = rh.Write(conn)         // write the response
		return true
	} else {
		_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
		return false
	}

}
