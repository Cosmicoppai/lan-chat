package chat

import (
	"crypto/sha1"
	"encoding/base64"
)

const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func computeAcceptKey(key string) string {
	concString := key + magicString
	hashedString := sha1.New()
	hashedString.Write([]byte(concString))
	return base64.StdEncoding.EncodeToString(hashedString.Sum(nil))
}
