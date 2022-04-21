package users

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"lan-chat/admin"
	"lan-chat/admin/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

var Tokens []string

func Handler(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	switch uri {
	case "/users":
		if r.Method == http.MethodGet {
			listUsers(w, r)
		}
	case "/user":
		switch r.Method {
		case http.MethodPost:
			registerUser(w, r)
		case http.MethodGet:
			listUser(w, r)
		case http.MethodPut:
			updateUsername(w, r)
		case http.MethodDelete:
			deleteUser(w, r)

		}

	}

}

func registerUser(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Data is in Invalid Format", http.StatusUnprocessableEntity)
		return
	}
	hashPassword := hashPass(user.Password)
	_, err = admin.Db.Exec("INSERT INTO lan_show.users VALUES ($1, $2)", user.Username, hashPassword)
	if err != nil {
		log.Println("Error while creating user", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {

	cred := r.Header.Get("Authorization")
	if cred != "" && strings.HasPrefix(cred, "Basic") {
		cred = strings.TrimPrefix(cred, "Basic")
		_decodedCred, err := base64.StdEncoding.DecodeString(cred)
		if err != nil {
			log.Println("Error while decoding authorization header", err)
			http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
			return
		}
		decodedCred := string(_decodedCred)
		sepIndex := strings.Index(decodedCred, ":")
		username, pass := decodedCred[:sepIndex], decodedCred[sepIndex+1:]
		if err = checkCredentials(w, username, pass); err != nil {
			token := getToken(username)
			tokenData := map[string]string{"token": token}
			_ = json.NewEncoder(w).Encode(tokenData)
			// _, _ = w.Write([]byte(token))
			return
		}
	}
	http.Error(w, "Authorization Header not present", http.StatusBadRequest)

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
		username, err := jwt.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		res, err := admin.Db.Exec("DELETE FROM lan_show.users where username=$1", username)
		if err != nil {
			log.Println("Error while deleting user", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)

}

func updateUsername(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token != "" {
		token = strings.TrimPrefix(token, "Bearer ")
		username, err := jwt.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		data := make(map[string]interface{})
		err = json.NewDecoder(r.Body).Decode(&data)
		if newUsername, ok := data["username"].(string); ok && err == nil {
			rows, err := admin.Db.Exec("UPDATE lan_show.users SET username=$1 WHERE username=$2", username, newUsername)
			if rowsAffected, _ := rows.RowsAffected(); rowsAffected > 1 && err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)

		} else {
			http.Error(w, "Invalid data format", http.StatusUnprocessableEntity)
		}
	}
}

func getToken(username string) string {

	claimsMap := map[string]string{
		"sub": username,
		"iss": "lan-chat",
		"exp": fmt.Sprint(time.Now().Add(time.Hour)),
	}

	header := "HS256"
	token, err := jwt.GenerateToken(header, claimsMap)
	if err != nil {
		log.Println("Error while generating token", err)
		return token
	}
	Tokens = append(Tokens, token)
	return token
}

func listUsers(w http.ResponseWriter, r *http.Request) {

}

func listUser(w http.ResponseWriter, r *http.Request) {

}

func checkCredentials(w http.ResponseWriter, username string, pass string) error {
	hashedPass := hashPass(pass)
	row, err := admin.Db.Query("SELECT password FROM lan_show.users WHERE username=$1", username)
	if err != nil {
		log.Println("Error in extracting hashed password: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return err
	}
	_ = row.Scan(&pass)
	if pass == hashedPass {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	http.Error(w, "", http.StatusUnauthorized)
	return fmt.Errorf("unauthorized")

}

func hashPass(p string) string {
	p = p + admin.Secret
	return fmt.Sprintf("%x", sha256.Sum256([]byte(p)))
}
