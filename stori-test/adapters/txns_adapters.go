package adapters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ortaman/stori-test/entity"
)

const (
	inputSeparator  = "/"
	outputSeparator = "-"
)

func validateTransaction(csvLine []string) (entity.Transaction, error) {
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

func ValidateTransactionsFromCSV(csvLines *[][]string) ([]entity.Transaction, error) {
	var transactionData []entity.Transaction

	for index, csvLine := range *csvLines {

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
