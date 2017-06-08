package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

//roleUser is used for assigning role to an user
const roleUser = "user"

var db *sql.DB

func init() {
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSslMode := "disable"
	if dbUser == "" {
		dbUser = "cubicasa"
	}
	if dbName == "" {
		dbName = "conversionV1"
	}
	if dbHost == "" {
		dbHost = "conversionv1.cwnxsiqll2vc.us-west-2.rds.amazonaws.com:5432"
	}
	if dbPassword == "" {
		dbPassword = "20Cubicasa16"
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
