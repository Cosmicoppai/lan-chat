package chat

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
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

	// buff := make([]byte, 4096)      // make buffer of 1024 bytes
	if connectionUpgrade(conn) { // if connection upgrade is successful
		read := bufio.NewReaderSize(conn, 5012)

		for {
			byte1, err := read.ReadByte()
			if err != nil {
				if err == io.EOF { // if user has closed the connection
					break
				}
				fmt.Println("Error!: ", err)
				return
			}

			isFinalBit := byte1 & finalBit
			opCode := byte1 & 0xf
			byte2, _ := read.ReadByte()        // Read the next byte
			isMaskBitSet := byte2&maskBit != 0 // check if first bit of 2nd byte is set or not

			var payloadLen int
			var nbs []byte // series of bytes to store if payLoadLen > 125
			if int(byte2&0x7f) <= 125 {
				payloadLen = int(byte2 & 0x7f)
			} else if byte2&0x7f == 126 { // Read the next 16 bits
				nb1, _ := read.ReadByte()
				nb2, _ := read.ReadByte()
				nbs = append(nbs, []byte{nb1, nb2}...)
				payloadLen = int(binary.BigEndian.Uint16(nbs))
			} else {
				for i := 0; i < 8; i++ { // Read the next 64 bits
					_nb, _ := read.ReadByte()
					nbs = append(nbs, _nb)
				}
				payloadLen = int(binary.BigEndian.Uint64(nbs))
			}

			if isMaskBitSet && (opCode == TextMessage) && isFinalBit != 0 {
				const maskKeyLen int = 4 // length of maskKeyLen
				remBytes, _ := read.Peek(maskKeyLen + payloadLen)
				maskKey := remBytes[:maskKeyLen]
				encodedPayLoad := remBytes[4:]
				decodedPayload := make([]byte, payloadLen)

				for i := 0; i < payloadLen; i++ {
					decodedPayload[i] = encodedPayLoad[i] ^ maskKey[i%maskKeyLen]

				}
				_, _ = read.Discard(read.Buffered()) // Reset the buffer

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
				_, err = conn.Write(finMsg)
				if err != nil {
					fmt.Println(err)
				}

			} else {
				_, _ = read.Discard(read.Buffered())
			}
		}
		conn.Close()
	}
	// t := http.ResponseWriter(conn)
	_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
}
