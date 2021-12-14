package chat

import (
	"net/http"
	"strings"
)

func ConnectionUpgrade(w http.ResponseWriter, r *http.Request) {

	if strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket" && r.Method == "GET" {

		webSockKey := r.Header.Get("Sec-WebSocket-Key")
		if webSockKey != "" {
			w.WriteHeader(http.StatusSwitchingProtocols)
			w.Header().Set("Upgrade", "websocket")
			w.Header().Set("Connection", "Upgrade")
			w.Header().Set("Sec-WebSocket-Version", "13")

			// The Sec-WebSocket-Accept header is important in that the server must derive it from the sec-WebSocket-Key
			// that the client sent to it

			secAcceptKey := computeAcceptKey(webSockKey) // compute the sec-WebSocket-Accept Key
			w.Header().Set("Sec-WebSocket-Accept", secAcceptKey)
			w.Header().Set("Sec-WebSocket-Origin", "http://192.168.1.106/")
			w.Header().Set("Sec-Websocket-Location", "ws://192.168.1.106/start-chat")
			// w.Write([]byte("Ok"))

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest) // if upgrade is not requested
	}

}
