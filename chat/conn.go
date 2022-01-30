package chat

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func Serve(conn net.Listener) error {
	for {
		req, e := conn.Accept()
		if e != nil {
			log.Println("Error while Accepting request", e)
			continue
		}
		go startChat(req)
	}
}

func startChat(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	if connectionUpgrade(conn) { // if connection upgrade is successful
		reader := bufio.NewReaderSize(conn, 5000012)
		var streamBuffer []byte

		for {
			byte1, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF || err.Error() == fmt.Sprintf("read tcp %s->%s: use of closed network connection", conn.LocalAddr().String(), conn.RemoteAddr().String()) { // if user has closed the connection
					fmt.Println(conn.RemoteAddr().String(), "has closed the connection")
					break
				}
				fmt.Println("Error!: ", err)
				continue
			}

			isFinalBit := byte1 & finalBit
			opCode := byte1 & 0xf
			byte2, _ := reader.ReadByte()      // Read the next byte
			isMaskBitSet := byte2&maskBit != 0 // check if first bit of 2nd byte is set or not

			if !isMaskBitSet {
				log.Println("mask bit is not set")
				_ = conn.Close() // if maskBit is not set, close the connection
			}

			var payloadLen int
			var nbs []byte // series of bytes to store if payLoadLen > 125
			if int(byte2&0x7f) <= 125 {
				payloadLen = int(byte2 & 0x7f)
			} else if byte2&0x7f == 126 { // Read the next 16 bits
				nb1, _ := reader.ReadByte() // Read 3rd byte
				nb2, _ := reader.ReadByte() // Read 4th byte
				nbs = append(nbs, []byte{nb1, nb2}...)
				payloadLen = int(binary.BigEndian.Uint16(nbs))
			} else {
				for i := 0; i < 8; i++ { // Read the next 64 bits(8 byte)
					_nb, _ := reader.ReadByte()
					nbs = append(nbs, _nb)
				}
				payloadLen = int(binary.BigEndian.Uint64(nbs))
			}

			decodedPayload := make([]byte, payloadLen)
			remBytes := make([]byte, maskKeyLen+payloadLen) // create a buffer of length(maskKeyLen + payloadLen)

			if opCode == TextMessage || opCode == continuationFrame || opCode == PingMessage {
				_, _ = io.ReadFull(reader, remBytes)
				maskKey := remBytes[:maskKeyLen]
				encodedPayLoad := remBytes[maskKeyLen:]
				for i := 0; i < payloadLen; i++ {
					decodedPayload[i] = encodedPayLoad[i] ^ maskKey[i%maskKeyLen]
				}
				if opCode == PingMessage {
					msg := []byte{finalBit | PongMessage}
					finMsg := append(msg, decodedPayload...)
					_, _ = conn.Write(finMsg)

				} else if (opCode == TextMessage) && isFinalBit != 0 { // for text msgs if finalBit is set, process the message
					handleChat(conn, decodedPayload, isFinalBit, opCode)
				} else if opCode == continuationFrame || isFinalBit == continuationFrame { // if opCode or finalBit is 0
					streamBuffer = append(streamBuffer, decodedPayload...) // append the msg
					if isFinalBit != 0 {                                   // when finalBit is Set
						handleChat(conn, streamBuffer, isFinalBit, 1)
						streamBuffer = nil // clear the stream Buffer
					}
				}

			} else {
				_, _ = reader.Discard(reader.Buffered())
			}
		}
		if _, exist := Data[conn]; exist {
			deleteUser(conn, UserRequest{}, false)
		}
	}
	_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
	log.Println(conn.RemoteAddr().String(), ": upgrade not successful")
}
