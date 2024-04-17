package usecase

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ortaman/stori-test/entity"
)

func validateTransaction(csvLine []string) (entity.Transaction, error) {
	const (
		inputSeparator  = "/"
		outputSeparator = "-"
	)

	// validate ID
	id, err := strconv.Atoi(csvLine[0])

	if err != nil {
		return entity.Transaction{}, fmt.Errorf("%s", err.Error())
	}

	// validate transaction amount
	transactionAmount, err := strconv.ParseFloat(csvLine[2], 64)

	if err != nil {
		return entity.Transaction{}, fmt.Errorf("%s", err.Error())
	}

	// validate date of the transaction
	dateSplited := strings.Split(csvLine[1], inputSeparator)
	dateSplited = append(dateSplited, strconv.Itoa(time.Now().Year()))

	for index, str := range dateSplited {
		if _, err := strconv.Atoi(dateSplited[index]); err != nil {
			return entity.Transaction{}, fmt.Errorf("%s", err.Error())
		}

		if len(str) == 1 {
			dateSplited[index] = "0" + str
		}
	}

	dateString := dateSplited[2] + outputSeparator + dateSplited[0] + outputSeparator + dateSplited[1]

	date, _ := time.Parse(time.DateOnly, dateString)

	// Return transaction struct parsed
	transaction := entity.Transaction{
		ID:          id,
		Date:        date,
		Transaction: transactionAmount,
	}
	return transaction, nil
}

func LoadTransactions(csvDir string) ([]entity.Transaction, error) {
	var transactionData []entity.Transaction

	csvFile, err := os.Open(csvDir)

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	csvLines, err := reader.ReadAll()

	if err != nil {
		return []entity.Transaction{}, fmt.Errorf("%s", err.Error())
	}

	for index, csvLine := range csvLines {

		if index != 0 {
			transaction, err := validateTransaction(csvLine)

			if err != nil {
				return []entity.Transaction{}, fmt.Errorf("%s", err.Error())
			}

			transactionData = append(transactionData, transaction)
		}
	}

	return transactionData, nil

}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetTransactionsResume(transactions []entity.Transaction) (entity.TransactionResume, error) {
	var (
		averageDebit    float64 = 0.0
		averageCredit   float64 = 0.0
		totalBalance    float64 = 0.0
		txsNumberDebit  float64 = 0.0
		txsNumberCredit float64 = 0.0
	)
	txnsNumberPerMonth := make(map[string]int)

	for _, txns := range transactions {

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
	from := "ente011@gmail.com"
	password := "mjjr aekr oikv qclv"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

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
