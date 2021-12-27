package chat

import (
	"bufio"
	"fmt"
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

	// buff := make([]byte, 1024)      // make buffer of 1024 bytes
	if connectionUpgrade(conn) { // if connection upgrade is successful
		read := bufio.NewReader(conn)

		for {
			peek, err := read.ReadByte()
			if err != nil {
				fmt.Println(err)
				return
			}
			// fmt.Println("byte1", peek)
			isFinalBit := peek & finalBit
			opCode := peek & 0xf
			peek2, _ := read.ReadByte()        // Read the next byte
			isMaskBitSet := peek2&maskBit != 0 // check if first bit of 2nd byte is set or not
			payloadLen := int(peek2 & 0x7f)
			// fmt.Println("finalBit", isFinalBit, "isMaskBitSet",isMaskBitSet, "opCode", opCode, "payLen", payloadLen)
			if isMaskBitSet && (opCode == TextMessage) && isFinalBit != 0 {
				const maskKeyLen int = 4 // length of maskKeyLen
				remBytes, _ := read.Peek(maskKeyLen + payloadLen)
				maskKey := remBytes[:maskKeyLen]
				encodedPayLoad := remBytes[4:]
				decodedPayload := make([]byte, payloadLen)

				for i := 0; i < payloadLen; i++ {
					decodedPayload[i] = encodedPayLoad[i] ^ maskKey[i%maskKeyLen]

				}
				// fmt.Println(string(decodedPayload))
				_, _ = read.Discard(read.Buffered()) // Reset the buffer
				msg := []byte{isFinalBit | opCode, byte(payloadLen)}
				finMsg := append(msg, decodedPayload...)
				_, _ = conn.Write(finMsg)

			} else {
				_, _ = read.Discard(read.Buffered())
			}
		}
	}
	// t := http.ResponseWriter(conn)
	_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
}
