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
	defer conn.Close()

	// buff := make([]byte, 4096)      // make buffer of 1024 bytes
	if connectionUpgrade(conn) { // if connection upgrade is successful
		read := bufio.NewReaderSize(conn, 5000012)

		for {
			byte1, err := read.ReadByte()
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
			byte2, _ := read.ReadByte()        // Read the next byte
			isMaskBitSet := byte2&maskBit != 0 // check if first bit of 2nd byte is set or not

			if !isMaskBitSet {
				log.Println("mask bit is not set")
				conn.Close()
			}

			var payloadLen int
			var nbs []byte // series of bytes to store if payLoadLen > 125
			if int(byte2&0x7f) <= 125 {
				payloadLen = int(byte2 & 0x7f)
			} else if byte2&0x7f == 126 { // Read the next 16 bits
				nb1, _ := read.ReadByte() // Read 3rd byte
				nb2, _ := read.ReadByte() // Read 4th byte
				nbs = append(nbs, []byte{nb1, nb2}...)
				payloadLen = int(binary.BigEndian.Uint16(nbs))
			} else {
				for i := 0; i < 8; i++ { // Read the next 64 bits(8 byte)
					_nb, _ := read.ReadByte()
					nbs = append(nbs, _nb)
				}
				payloadLen = int(binary.BigEndian.Uint64(nbs))
			}

			if isMaskBitSet && (opCode == TextMessage) && isFinalBit != 0 {
				remBytes, _ := read.Peek(maskKeyLen + payloadLen)
				maskKey := remBytes[:maskKeyLen]
				encodedPayLoad := remBytes[4:]
				decodedPayload := make([]byte, payloadLen)
				for i := 0; i < payloadLen; i++ {
					decodedPayload[i] = encodedPayLoad[i] ^ maskKey[i%maskKeyLen]

				}
				handleChat(conn, decodedPayload, isFinalBit, opCode)
				_, _ = read.Discard(read.Buffered()) // Reset the buffer

			} else {
				_, _ = read.Discard(read.Buffered())
			}
		}
		deleteUser(conn, UserRequest{}, false)
	}
	// t := http.ResponseWriter(conn)
	_, _ = conn.Write([]byte(string(rune(http.StatusBadRequest))))
	log.Fatalln("upgrade not successful")
}
