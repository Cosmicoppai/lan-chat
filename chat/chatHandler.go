package chat

import (
	"bufio"
	"encoding/binary"
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
		read := bufio.NewReaderSize(conn, 40966)

		for {
			byte1, err := read.ReadByte()
			if err != nil {
				fmt.Println("Error!: ", err)
				return
			}

			isFinalBit := byte1 & finalBit
			opCode := byte1 & 0xf
			byte2, _ := read.ReadByte()        // Read the next byte
			isMaskBitSet := byte2&maskBit != 0 // check if first bit of 2nd byte is set or not

			var payloadLen int
			if int(byte2&0x7f) <= 125 {
				payloadLen = int(byte2 & 0x7f)
			} else if byte2&0x7f == 126 { // Read the next 16 bits
				nb1, _ := read.ReadByte()
				nb2, _ := read.ReadByte()
				mbs := []byte{nb1, nb2}
				payloadLen = int(binary.BigEndian.Uint16(mbs))
			} else {
				var nb []byte
				for i := 0; i < 8; i++ { // Read the next 64 bits
					_nb, _ := read.ReadByte()
					nb = append(nb, _nb)
				}
				payloadLen = int(binary.BigEndian.Uint64(nb))
			}

			if isMaskBitSet && (opCode == TextMessage) && isFinalBit != 0 {
				fmt.Println("payLoad", payloadLen)
				const maskKeyLen int = 4 // length of maskKeyLen
				remBytes, _ := read.Peek(maskKeyLen + payloadLen)
				maskKey := remBytes[:maskKeyLen]
				encodedPayLoad := remBytes[4:]
				decodedPayload := make([]byte, payloadLen)

				for i := 0; i < payloadLen; i++ {
					decodedPayload[i] = encodedPayLoad[i] ^ maskKey[i%maskKeyLen]

				}
				fmt.Println(string(decodedPayload))
				_, _ = read.Discard(read.Buffered()) // Reset the buffer
				msg := []byte{isFinalBit | opCode, byte(payloadLen)}
				finMsg := append(msg, decodedPayload...)
				_, _ = conn.Write(finMsg)
				fmt.Println(read)

			} else {
				_, _ = read.Discard(read.Buffered())
			}
		}
	}
	// t := http.ResponseWriter(conn)
	_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
}
