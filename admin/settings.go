package admin

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
)

var (
	hostname     string
	port         string
	username     string
	password     string
	databaseName string
	Secret       string
	Db           *sql.DB
)

func setConfig() {
	hostname = os.Getenv("hostname")
	port = os.Getenv("port")
	username = os.Getenv("username")
	password = os.Getenv("password")
	databaseName = os.Getenv("databaseName")
	Secret = os.Getenv("secret")
}

func InitializeDB() {
	var err error
	setConfig()
	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		hostname, port, username, password, databaseName)
	Db, err = sql.Open("postgres", pgConnString)
	if err != nil {
		log.Println("Failed to connect to the database")
		log.Fatalln(err)
	}
	log.Println("Successfully Connected with database: ")
	query, err := ioutil.ReadFile("admin/model.sql")
	if err != nil {
		log.Println("Error while reading sql script: ")
		log.Fatalln(err)
	}
	if _, err := Db.Exec(string(query)); err != nil {
		log.Println("Error while Executing the Script: ")
		log.Fatalln(err)
	}

}
