package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"lan-chat/admin"
	"lan-chat/logger"
	"strings"
	"time"
)

var secret = admin.Secret

var Tokens = make(map[string][]string)

type Claims struct {
	Sub     string      `json:"sub"`
	IsAdmin bool        `json:"isAdmin"`
	Exp     json.Number `json:"exp"`
	ISS     string      `json:"iss"`
}

func tokenPresent(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

var InvalidToken = errors.New("invalid token")

// GenerateToken : This will create New Token
func GenerateToken(header map[string]string, payload map[string]interface{}) (string, error) {
	// create a new hash of type sha256. We pass the secret key to it
	h := hmac.New(sha256.New, []byte(secret))
	headerStr, err := json.Marshal(header)
	if err != nil {
		fmt.Println("Error generating Token")
		return "", err
	}
	header64 := base64.StdEncoding.EncodeToString(headerStr)

	payloadStr, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error generating Token")
		return "", err
	}
	payload64 := base64.StdEncoding.EncodeToString(payloadStr)

	// add the encoded string.
	message := header64 + "." + payload64

	// We have the unsigned message ready.
	unsignedStr := string(headerStr) + string(payloadStr)

	// write this to the SHA256 to hash it.
	h.Write([]byte(unsignedStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	tokenStr := message + "." + signature
	username := payload["sub"].(string)

	Tokens[username] = append(Tokens[username], tokenStr)
	return tokenStr, nil
}

// ValidateToken : This helps in validating the token
func ValidateToken(token string) (Claims, error) {
	claims := Claims{}
	// JWT has 3 parts separated by '.'
	splitToken := strings.Split(token, ".")
	// if length is not 3, we know that the token is corrupt
	if len(splitToken) != 3 {
		return claims, InvalidToken
	}

	// decode the header and payload back to strings
	header, err := base64.StdEncoding.DecodeString(splitToken[0])
	if err != nil {
		return claims, InvalidToken
	}
	payload, err := base64.StdEncoding.DecodeString(splitToken[1])
	if err != nil {
		return claims, InvalidToken
	}
	//again create the signature
	unsignedStr := string(header) + string(payload)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedStr))

	err = json.Unmarshal(payload, &claims) // deserialize the payloadData
	if err != nil {
		logger.ErrorLog.Println("Error in unmarshalling the payload: ", err)
		return claims, err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// check signature, expire time and token existence
	expTime, _ := claims.Exp.Int64()

	username := claims.Sub
	if signature != splitToken[2] || !tokenPresent(token, Tokens[username]) || expTime < time.Now().Unix() {
		return claims, InvalidToken
	}

	return claims, nil
}

func DeleteUserTokens(username string) {
	delete(Tokens, username)
}
