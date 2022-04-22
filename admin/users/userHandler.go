package users

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"lan-chat/admin"
	"lan-chat/admin/dbErrors"
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
	if err != nil || user.Username == "" || user.Password == "" {
		log.Println(err)
		http.Error(w, "Data is in Invalid Format", http.StatusUnprocessableEntity)
		return
	}
	hashedPassword := hashPass(user.Password)
	err = InsertUser(user.Username, hashedPassword)
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Class() == "23" { // if error is about integrity constraint violation
			http.Error(w, "username taken", http.StatusConflict)
			return
		}
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
		if checkCredentials(w, username, pass) {
			token := getToken(username)
			tokenData := map[string]string{"token": token}
			_ = json.NewEncoder(w).Encode(tokenData)
			return
		}
		return
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
		_, err = admin.Db.Exec("DELETE FROM lan_show.users where username=$1", username)
		if dbErrors.InternalServerError(err) {
			log.Println("Error while deleting user", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Unauthorized", http.StatusUnauthorized)

}

func updateUsername(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(*User).Username
	log.Println(username)
	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
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
	usersList := map[string][]string{"users": {}}
	var user string
	rows, err := admin.Db.Query("SELECT username FROM lan_show.users;")
	if dbErrors.InternalServerError(err) {
		log.Println("Error in extracting hashed password: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user)
		if err != nil {
			log.Println(err)
		}
		usersList["users"] = append(usersList["users"], user)
	}
	err = json.NewEncoder(w).Encode(usersList)
	if err != nil {
		log.Println("Error while encoding the data into Json: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func listUser(w http.ResponseWriter, r *http.Request) {
	var user string
	row := admin.Db.QueryRow("SELECT username FROM lan_show.users;")
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No user Exist", http.StatusNotFound)
			return
		}
		log.Println(err)
	}
	err = json.NewEncoder(w).Encode(map[string]string{"user": user})
	if err != nil {
		log.Println("Error while encoding the data into Json: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func checkCredentials(w http.ResponseWriter, username string, pass string) bool {
	hashedPass := hashPass(pass)
	row, err := admin.Db.Query("SELECT password FROM lan_show.users WHERE username=$1", username)
	if dbErrors.InternalServerError(err) {
		log.Println("Error in extracting hashed password: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return false
	}
	_ = row.Scan(&pass)
	if pass == hashedPass {
		w.WriteHeader(http.StatusOK)
		return true
	}
	http.Error(w, "", http.StatusUnauthorized)
	return false

}

func hashPass(p string) string {
	p = p + admin.Secret
	return fmt.Sprintf("%x", sha256.Sum256([]byte(p)))
}

func CreateSuperUser(username string, password string) {
	hashedPassword := hashPass(password)
	err := InsertUser(username, hashedPassword)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Super User Successfully created ....")
}

func InsertUser(username string, password string) error {
	log.Println(username, password)
	_, err := admin.Db.Exec("INSERT INTO lan_show.users VALUES ($1, $2)", username, password)
	return err
}
