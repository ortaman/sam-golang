package entity

import "time"

type (
	Transaction struct {
		ID          int
		Date        time.Time
		Transaction float64
	}
	TransactionResume struct {
		TotalBalance       float64
		AverageDebit       float64
		AverageCredit      float64
		TxnsNumberPerMonth map[string]int
	}
)

type TnxsUseCaseI interface {
	SaveTransactions(email string, transactions *[]Transaction) error
	SendEmail(transactions *[]Transaction, emails_to []string, termplateDir string) error
}

type TnxsRepoI interface {
	GetUserByEmail(email string) int
	CreateUser(email string) int
	SaveTransactions(customer_id int, transactions *[]Transaction) error
}

type (
	Email struct {
	}
)

type EmailRepoI interface {
	SendEmail(transactionsResume *TransactionResume, emails_to []string, termplateDir string) error
}
