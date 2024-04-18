package usecase

import (
	"bytes"
	"fmt"
	"math"
	"net/smtp"
	"os"
	"text/template"

	"github.com/ortaman/stori-test/entity"
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

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
		TotalBalance:       roundFloat(totalBalance, 2),
		AverageDebit:       roundFloat(averageDebit, 2),
		AverageCredit:      roundFloat(averageCredit, 2),
		TxnsNumberPerMonth: txnsNumberPerMonth,
	}, nil
}

func SendEmail(email_to []string, termplateDir string, transactionResume *entity.TransactionResume) error {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_CODE")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Message
	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Stori Test \n%s\n\n", mimeHeaders)))

	template, TemplatError := template.ParseFiles(termplateDir)

	if TemplatError != nil {
		return fmt.Errorf("%s", TemplatError.Error())
	}

	template.Execute(&body, transactionResume)

	// Authentication and send email
	auth := smtp.PlainAuth("", from, password, smtpHost)
	emailError := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, email_to, body.Bytes())

	if emailError != nil {
		return fmt.Errorf("%s", emailError.Error())
	}

	return nil
}
