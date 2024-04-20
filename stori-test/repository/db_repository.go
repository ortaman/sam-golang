package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ortaman/stori-test/entity"
)

type DBRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) entity.UserRepository {
	return &DBRepository{db}
}

func (dbRepository *DBRepository) CreateUser(email string) int {

	sqlStatement := `INSERT INTO customer (email) VALUES ($1) RETURNING customer_id`

	var userId int
	err := dbRepository.db.QueryRow(sqlStatement, email).Scan(&userId)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	return userId
}

func (dbRepository *DBRepository) GetUserByEmail(email string) int {

	sqlStatement := `SELECT customer_id FROM customer WHERE email=($1)`

	var userId int
	err := dbRepository.db.QueryRow(sqlStatement, email).Scan(&userId)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return -1
	case nil:
		return userId
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return userId
}
