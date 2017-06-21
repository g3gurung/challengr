package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/challengr/middleware"
	_ "github.com/lib/pq"
)

//roleUser is used for assigning role to an user
const roleUser = "user"

//JWTSecret is user for encrypting and decrypting jwt
const JWTSecret = middleware.JWTSecret

var db *sql.DB

func init() {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSslMode := "disable"
	if dbUser == "" {
		dbUser = "challengr"
	}
	if dbName == "" {
		dbName = "challengrdb"
	}
	if dbHost == "" {
		dbHost = "challengr.cd61ijduodvj.us-east-1.rds.amazonaws.com:5432"
	}
	if dbPassword == "" {
		dbPassword = "20Challengr17"
	}

	dataSource := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbName, dbSslMode)

	log.Printf("dataSource: %v", dataSource)

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Println("DB ping started...")
	if err = db.Ping(); err != nil {
		log.Printf("DB ping failed with error...%v", err)
	}
	log.Println("DB connected.")
}
