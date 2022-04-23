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
	"lan-chat/httpErrors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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

func registerUser(w http.ResponseWriter, r *http.Request) { // only admin can register a user
	claims := r.Context().Value("claims").(jwt.Claims)
	if !claims.IsAdmin {
		httpErrors.Forbidden(w)
		return
	}
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		log.Println(err)
		//http.Error(w, "Data is in Invalid Format", http.StatusUnprocessableEntity)
		httpErrors.UnProcessableEntry(w, "Data is in Invalid Format")
		return
	}
	hashedPassword := hashPass(user.Password)
	err = insertUser(user.Username, hashedPassword)
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

func Login(w http.ResponseWriter, r *http.Request) { // anyone with their account can Log In

	if r.Method != http.MethodPost {
		httpErrors.MethodNotAllowed(w)
		return
	}

	cred := r.Header.Get("Authorization")
	if cred != "" && strings.HasPrefix(cred, "Basic") {
		cred = strings.TrimPrefix(cred, "Basic ")
		_decodedCred, err := base64.StdEncoding.DecodeString(cred)
		if err != nil {
			log.Println("Error while decoding authorization header", err)
			http.Error(w, "Invalid Authorization Header", http.StatusBadRequest)
			return
		}
		decodedCred := string(_decodedCred)
		sepIndex := strings.Index(decodedCred, ":")
		username, pass := decodedCred[:sepIndex], decodedCred[sepIndex+1:]
		if isAdmin, credOk := checkCredentials(w, username, pass); credOk {
			token := getToken(username, isAdmin)
			tokenData := map[string]string{"token": token}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(tokenData)
			return
		}
		return
	}
	http.Error(w, "Authorization Header not present", http.StatusBadRequest)

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("claims").(jwt.Claims).Sub
	_, err := admin.Db.Exec("DELETE FROM lan_show.users where username=$1", username)
	if dbErrors.InternalServerError(err) {
		log.Println("Error while deleting user", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	jwt.DeleteUserTokens(username)
	w.WriteHeader(http.StatusOK)
	return
}

func updateUsername(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("claims").(jwt.Claims).Sub
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

func ListUsers(w http.ResponseWriter, r *http.Request) {
	usersList := map[string][]string{"users": {}}
	var user string
	rows, err := admin.Db.Query("SELECT username FROM lan_show.users;")
	if dbErrors.InternalServerError(err) {
		log.Println("Error in extracting users:  ", err)
		httpErrors.InternalServerError(w)
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
		httpErrors.InternalServerError(w)
		return
	}

}

func listUser(w http.ResponseWriter, r *http.Request) {
	uri, _ := url.Parse(r.RequestURI)
	queryParams, _ := url.ParseQuery(uri.RawQuery)
	username := queryParams.Get("username")
	if username == "" {
		httpErrors.BadRequest(w)
		return
	}
	var user string
	row := admin.Db.QueryRow("SELECT username FROM lan_show.users where username=$1;", username)
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

func checkCredentials(w http.ResponseWriter, username string, pass string) (isAdmin bool, credValid bool) {
	isAdmin = false
	hashedPass := hashPass(pass)
	// var passFromDb string
	row := admin.Db.QueryRow("SELECT password, isAdmin FROM lan_show.users WHERE username=$1", username)

	err := row.Scan(&pass, &isAdmin)
	if dbErrors.InternalServerError(err) {
		log.Println("Error in extracting hashed password: ", err)
		http.Error(w, "", http.StatusInternalServerError)
		return isAdmin, false
	}
	if pass == hashedPass {
		return isAdmin, true
	}
	httpErrors.UnProcessableEntry(w)
	return isAdmin, false

}

func hashPass(p string) string {
	p = p + admin.Secret
	return fmt.Sprintf("%x", sha256.Sum256([]byte(p)))
}

func CreateSuperUser(username string, password string) {
	hashedPassword := hashPass(password)
	err := insertUser(username, hashedPassword, true)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Super User Successfully created ....")
}

func insertUser(username string, password string, _admin ...bool) error {
	isAdmin := false
	if len(_admin) == 1 {
		isAdmin = true
	}
	_, err := admin.Db.Exec("INSERT INTO lan_show.users VALUES ($1, $2, $3)", username, password, isAdmin)
	return err
}

func getToken(username string, isAdmin bool) string {

	claimsMap := map[string]interface{}{
		"sub": username,
		"iss": "lan-chat",
		"exp": fmt.Sprint(time.Now().Add(300 * time.Minute).Unix()),
	}
	if isAdmin {
		claimsMap["isAdmin"] = true
	}
	log.Println(claimsMap)

	header := map[string]string{"alg": "HS256"}
	token, err := jwt.GenerateToken(header, claimsMap)
	if err != nil {
		log.Println("Error while generating token", err)
		return token
	}
	return token
}
