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

	// psqlInfo := "host=localhost port=5432 user=db_user password=db_pass dbname=db_stori sslmode=disable"
	// export DB_HOST="host.docker.internal" DB_PORT="5432" DB_USER="db_user" DB_PASS="db_pass" DB_NAME="db_stori" sslmode="disable" DB_DRIVER="postgres"

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
