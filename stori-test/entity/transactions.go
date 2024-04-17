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
