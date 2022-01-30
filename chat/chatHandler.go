package chat

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

type UserRequest struct {
	Typ       string `json:"typ"`
	UserName  string `json:"userName"`
	Msg       string `json:"msg"`
	Token     string `json:"token"`
	TotalUser int    `json:"totalUser"`
}

type UserData struct {
	userName string
	token    string
}

var Data = map[net.Conn]UserData{}

func handleChat(conn net.Conn, decodedPayload []byte, isFinalBit byte, opCode byte) {

	var req UserRequest
	decodeErr := json.Unmarshal(decodedPayload, &req)
	if decodeErr != nil {
		log.Println("Error while decoding data: ", decodeErr)
	}
	if strings.ToLower(req.Typ) == "add" {
		addUser(conn, req)
	} else if strings.ToLower(req.Typ) == "remove" {
		deleteUser(conn, req, true)
	} else if strings.ToLower(req.Typ) == "txt-msg" || strings.ToLower(req.Typ) == "img-msg" {
		_username := Data[conn].userName
		req.UserName = _username
		finMsg := encodeMsg(isFinalBit, opCode, req)
		sendMsg(finMsg)
	}

}

func addUser(conn net.Conn, request UserRequest) {
	var udata UserData
	username := request.UserName // userName from request
	if username != "" {
		if _, exist := Data[conn]; exist {
			_ = conn.Close()
		}
		if !checkUserName(username) {
			udata.userName = username
			udata.token = createToken() // set token

			_addUser(conn, udata)
			fmt.Printf("%s has joined the chat\n", username)
		} else {
			sendSingleMsg(conn, UserRequest{Typ: "error", Msg: "User-name already exists !"})
		}
	}

}

func _addUser(conn net.Conn, udata UserData) {

	// send msg to all connections
	msg := UserRequest{Typ: "alert",
		Msg: fmt.Sprintf("%s has joined the chat", udata.userName), TotalUser: len(Data) + 1}
	encodedMsg := encodeMsg(finalBit, TextMessage, msg)
	sendMsg(encodedMsg)
	Data[conn] = udata

	// send msg to the requested-user with token
	msgWithToken := UserRequest{Typ: "alert", Msg: fmt.Sprintf("%s has joined the chat", udata.userName),
		Token:     udata.token,
		TotalUser: len(Data)}
	sendSingleMsg(conn, msgWithToken)

}

func checkUserName(userName string) bool {
	for _, user := range Data {
		if user.userName == userName {
			return true
		}
	}
	return false
}

func deleteUser(conn net.Conn, request UserRequest, tokenRequired bool) {
	if !tokenRequired {
		_deleteUser(conn)
	} else if request.Token == Data[conn].token {
		sendSingleMsg(conn, UserRequest{Typ: "success", Msg: "Successfully left the chat"})
		_deleteUser(conn)
	} else {
		sendSingleMsg(conn, UserRequest{Typ: "error", Msg: "Invalid Token"})
	}
}

func _deleteUser(conn net.Conn) {
	msg := UserRequest{Typ: "alert", Msg: fmt.Sprintf("%s has left the chat", Data[conn].userName), TotalUser: len(Data) - 1}
	encodedMsg := encodeMsg(finalBit, TextMessage, msg)
	delete(Data, conn)
	sendMsg(encodedMsg)
}

func sendMsg(msg []byte) {
	for conn := range Data {
		_, err := conn.Write(msg)
		if err != nil {
			log.Println("Error in sendMsg", err)
			deleteUser(conn, UserRequest{}, false)
		}
	}

}

func sendSingleMsg(conn net.Conn, msg UserRequest) {
	encodedMsg := encodeMsg(finalBit, TextMessage, msg)
	_, err := conn.Write(encodedMsg)
	if err != nil {
		log.Fatalln(err)
	}
}

func encodeMsg(isFinalBit byte, opCode byte, msg UserRequest) []byte {
	_msg := []byte{isFinalBit | opCode} // add the first byte consist of isFinalBit + OpCode
	encodedPayload, _err := json.Marshal(msg)
	if _err != nil {
		log.Fatalln(_err)
	}
	var payloadLenBytes []byte
	payloadLen := len(encodedPayload)
	if payloadLen < 126 {
		payloadLenBytes = []byte{byte(payloadLen)}
	} else if len(encodedPayload) <= 65535 {
		payloadLenBytes = []byte{byte(126)}
		_payloadLenBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(_payloadLenBytes, uint16(payloadLen))
		payloadLenBytes = append(payloadLenBytes, _payloadLenBytes...)
	} else {
		payloadLenBytes = []byte{byte(127)}
		_payloadLenBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(_payloadLenBytes, uint64(payloadLen))
		payloadLenBytes = append(payloadLenBytes, _payloadLenBytes...)
	}
	finMsg := append(_msg, payloadLenBytes...) // add the bytes consist of payloadLen
	finMsg = append(finMsg, encodedPayload...) // add the message
	return finMsg
}

func createToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 64)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:64]

}
