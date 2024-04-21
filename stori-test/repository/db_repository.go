package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ortaman/stori-test/entity"
)

type DBRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) entity.TnxsRepoI {
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
		log.Println("No rows were returned!")
		return -1
	case nil:
		return userId
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return userId
}

func (dbRepository *DBRepository) SaveTransactions(customer_id int, transactions *[]entity.Transaction) error {

	placeholders := []string{}
	tnxsValues := []interface{}{}

	for index, tnxs := range *transactions {
		placeholders = append(
			placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", index*4+1, index*4+2, index*4+3, index*4+4))

		tnxsValues = append(tnxsValues, tnxs.ID, tnxs.Date, tnxs.Transaction, customer_id)
	}

	sqlQuery := `
		INSERT INTO
			transactions(id, trans_date, amount, customer_id)
		VALUES %s ON CONFLICT (id) DO UPDATE SET
		trans_date = EXCLUDED.trans_date, amount = EXCLUDED.amount, customer_id = EXCLUDED.customer_id`

	sqlQuery = fmt.Sprintf(sqlQuery, strings.Join(placeholders, ","))

	tx, err := dbRepository.db.Begin()
	if err != nil {
		log.Fatalf("Unable to inicialize database. %v", err)
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	if _, err = tx.Exec(sqlQuery, tnxsValues...); err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Unable to commit the query. %v", err)
	}

	return nil

}
