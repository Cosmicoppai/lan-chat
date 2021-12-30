package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type UserRequest struct {
	Typ      string `json:"typ"`
	UserName string `json:"username"`
	Msg      string `json:"msg"`
	Token    string `json:"token"`
}

type UserData struct {
	userName string
	token    string
}

var Data map[net.Conn]UserData

func handleChat(conn net.Conn, decodedPayload []byte, byte2 byte, nbs []byte, isFinalBit byte, opCode byte, payloadLen int) {

	var req UserRequest
	jsonDecodeMsg := json.Unmarshal(decodedPayload, &req)
	if jsonDecodeMsg != nil {
		log.Fatalln(jsonDecodeMsg)
		return
	}
	if strings.ToLower(req.Typ) == "add" {
		addUser(conn, req)
	} else if strings.ToLower(req.Typ) == "delete" {
		deleteUser(conn, req)
	} else {

		var payloadLenBytes []byte
		if payloadLen == int(byte2&0x7f) {
			payloadLenBytes = []byte{byte(payloadLen)}
		} else {
			payloadLenBytes = []byte{byte2 & 0x7f}
			payloadLenBytes = append(payloadLenBytes, nbs...)
		}
		msg := []byte{isFinalBit | opCode}         // add the first byte consist of isFinalBit + OpCode
		finMsg := append(msg, payloadLenBytes...)  // add the bytes consist of payloadLen
		finMsg = append(finMsg, decodedPayload...) // add the message
		_, err := conn.Write(finMsg)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func addUser(conn net.Conn, request UserRequest) {
	var udata UserData
	username := request.UserName
	if user, exist := Data[conn]; exist {
		if request.UserName != user.userName {
			conn.Close()
			return
		}
	}
	udata.userName = username
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 64)
	rand.Read(b)
	udata.token = fmt.Sprintf("%x", b)[:64]
	msg := UserRequest{Typ: "alert", Msg: fmt.Sprintf("%s has joined the chat", username), Token: udata.token}
	encodedMsg := encodeMsg(finalBit&1, TextMessage, msg)
	sendChat(encodedMsg)

}

func deleteUser(conn net.Conn, request UserRequest) {
	if request.Token == Data[conn].token {
		delete(Data, conn)
	}

}

func sendChat(msg []byte) {

}

func encodeMsg(isFinalBit byte, opCode byte, msg UserRequest) []byte {
	return []byte("lol")
}
