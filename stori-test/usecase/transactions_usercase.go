package usecase

import (
	"github.com/ortaman/stori-test/entity"
	"github.com/ortaman/stori-test/utils"
)

type TnxsUsercase struct {
	TnxsRepoI  entity.TnxsRepoI
	emailRepoI entity.EmailRepoI
}

func NewTnxsUsecase(tnxsRepoI entity.TnxsRepoI, emailRepoI entity.EmailRepoI) entity.TnxsUseCaseI {

	return &TnxsUsercase{
		TnxsRepoI:  tnxsRepoI,
		emailRepoI: emailRepoI,
	}
}

func (tnxsUsercase *TnxsUsercase) SaveTransactions(email string, transactions *[]entity.Transaction) error {

	customer_id := tnxsUsercase.TnxsRepoI.GetUserByEmail(email)

	if customer_id < 0 {
		customer_id = tnxsUsercase.TnxsRepoI.CreateUser(email)
	}

	if err := tnxsUsercase.TnxsRepoI.SaveTransactions(customer_id, transactions); err != nil {
		return err
	}

	return nil
}

func (tnxsUsercase *TnxsUsercase) GetTransactionsResume(email string, transactions *[]entity.Transaction) error {

	customer_id := tnxsUsercase.TnxsRepoI.GetUserByEmail(email)

	if customer_id < 0 {
		customer_id = tnxsUsercase.TnxsRepoI.CreateUser(email)
	}

	if err := tnxsUsercase.TnxsRepoI.SaveTransactions(customer_id, transactions); err != nil {
		return err
	}

	return nil
}

func (tnxsUsercase *TnxsUsercase) SendEmail(transactions *[]entity.Transaction, emails_to []string, termplateDir string) error {
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

	transactionResume := entity.TransactionResume{
		TotalBalance:       utils.RoundFloat(totalBalance, 2),
		AverageDebit:       utils.RoundFloat(averageDebit, 2),
		AverageCredit:      utils.RoundFloat(averageCredit, 2),
		TxnsNumberPerMonth: txnsNumberPerMonth,
	}

	tnxsUsercase.emailRepoI.SendEmail(&transactionResume, emails_to, termplateDir)

	return nil
}
