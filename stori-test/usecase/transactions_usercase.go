package usecase

import (
	"github.com/ortaman/stori-test/entity"
	"github.com/ortaman/stori-test/utils"
)

func GetTransactionsResume(transactions *[]entity.Transaction) (entity.TransactionResume, error) {
	var (
		averageDebit    float64 = 0.0
		averageCredit   float64 = 0.0
		totalBalance    float64 = 0.0
		txsNumberDebit  float64 = 0.0
		txsNumberCredit float64 = 0.0
	)
	txnsNumberPerMonth := make(map[string]int)

	for _, txns := range *transactions {

		txnsAmount := txns.Transaction

		txnsNumberPerMonth[txns.Date.Month().String()] += 1

		if txnsAmount > 0 {
			txsNumberDebit += 1
			averageDebit += txnsAmount
		} else {
			txsNumberCredit += 1
			averageCredit += txnsAmount
		}

		totalBalance += txnsAmount
	}

	if txsNumberDebit > 0 {
		averageDebit = averageDebit / txsNumberDebit
	}
	if txsNumberCredit > 0 {
		averageCredit = averageCredit / txsNumberCredit
	}

	return entity.TransactionResume{
		TotalBalance:       utils.RoundFloat(totalBalance, 2),
		AverageDebit:       utils.RoundFloat(averageDebit, 2),
		AverageCredit:      utils.RoundFloat(averageCredit, 2),
		TxnsNumberPerMonth: txnsNumberPerMonth,
	}, nil
}
