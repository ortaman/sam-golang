package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func NewMyPSQLConnection() *sql.DB {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		driver   = os.Getenv("DB_DRIVER")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASS")
		dbName   = os.Getenv("DB_NAME")
	)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	fmt.Println("psqlInfo: ", psqlInfo)
	db, err := sql.Open(driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	return db
}

func NewPSQLConfig() *map[string]string {

	config := map[string]string{
		"from":     os.Getenv("EMAIL_FROM"),
		"password": os.Getenv("EMAIL_CODE"),
		"smtpHost": os.Getenv("SMTP_HOST"),
		"smtpPort": string(os.Getenv("SMTP_PORT")),
	}
	return &config
}
