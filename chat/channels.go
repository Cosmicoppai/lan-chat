package chat

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	jsonDecodeMsg := json.Unmarshal(decodedPayload, &req)
	if jsonDecodeMsg != nil {
		conn.Close()
		log.Fatalln(jsonDecodeMsg)
	}
	if strings.ToLower(req.Typ) == "add" {
		addUser(conn, req)
	} else if strings.ToLower(req.Typ) == "remove" {
		deleteUser(conn, req, true)
	} else {
		_username := Data[conn].userName
		if _username != "" {
			req.UserName = _username
			finMsg := encodeMsg(isFinalBit, opCode, req)
			_, err := conn.Write(finMsg)
			if err != nil {
				fmt.Println(err)
			} else {
				conn.Close()
			}
		}
	}

}

func addUser(conn net.Conn, request UserRequest) {
	var udata UserData
	username := request.UserName // userName from request
	if user, exist := Data[conn]; exist {
		if request.UserName != user.userName {
			conn.Close()
			return
		} else {
			_addUser(username, user.token)
			return

		}
	}
	udata.userName = username
	udata.token = createToken() // set token
	Data[conn] = udata

	_addUser(username, udata.token)
	fmt.Printf("%s has joined the chat", username)

}

func _addUser(username string, token string) {
	msg := UserRequest{Typ: "alert",
		Msg:   fmt.Sprintf("%s has joined the chat", username),
		Token: token, TotalUser: len(Data)}
	encodedMsg := encodeMsg(finalBit, TextMessage, msg)
	sendMsg(encodedMsg)

}

func deleteUser(conn net.Conn, request UserRequest, tokenRequired bool) {
	if !tokenRequired {
		_deleteUser(conn)
	} else {
		if request.Token == Data[conn].token {
			_deleteUser(conn)
		}
	}

}

func _deleteUser(conn net.Conn) {
	msg := UserRequest{Typ: "alert", Msg: fmt.Sprintf("%s has left the chat", Data[conn].userName), TotalUser: len(Data) - 1}
	encodedMsg := encodeMsg(finalBit, TextMessage, msg)
	delete(Data, conn)
	sendMsg(encodedMsg)
}

func sendMsg(msg []byte) {
	for conn, _ := range Data {
		_, err := conn.Write(msg)
		if err != nil {
			log.Fatalln(err)
		}
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
	} else if len(encodedPayload) <= math.MaxInt16 {
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
