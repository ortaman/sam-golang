package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ortaman/stori-test/adapters"
	"github.com/ortaman/stori-test/entity"
	"github.com/ortaman/stori-test/infra"
	"github.com/ortaman/stori-test/repository"
	"github.com/ortaman/stori-test/usecase"

	_ "github.com/lib/pq"
)

const (
	csvDir      = "./txns.csv"
	templateDir = "./txns_template.html"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Validate the body request
	bodyMap := map[string]string{}

	if err := json.Unmarshal([]byte(request.Body), &bodyMap); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(entity.ErrorBody.Error(), err.Error()),
			StatusCode: 400,
		}, nil
	}

	// Validate that email exists
	email, emailExist := bodyMap["email"]

	if !emailExist {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(entity.ErrEmail.Error(), "email"),
			StatusCode: 400,
		}, nil
	}

	// Read csv file
	csvLinesPointer, err := adapters.ReadAllCSV(csvDir)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(entity.ErrCSV.Error(), err.Error()),
			StatusCode: 400,
		}, nil
	}

	// Validate csv file and load data in a slice
	transactionData, err := adapters.ValidateTransactionsFromCSV(csvLinesPointer)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(entity.ErrTrxsValid.Error(), err.Error()),
			StatusCode: 500,
		}, nil
	}

	db := infra.NewMyPSQLConnection()
	dbRepository := repository.NewSQLRepository(db)

	PSQConfig := infra.NewPSQLConfig()
	emailRepository := repository.NewEmailRepository(PSQConfig)

	tnxsUsercase := usecase.NewTnxsUsecase(dbRepository, emailRepository)

	// Calculate transactions
	if err := tnxsUsercase.SaveTransactions(email, &transactionData); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(entity.ErrTrxnSave.Error(), err.Error()),
			StatusCode: 500,
		}, nil
	}

	// Calculate Transaction resume and send email
	if err := tnxsUsercase.SendEmail(&transactionData, []string{email}, templateDir); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf(entity.ErrTrxnSave.Error(), err.Error()),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Transactions Resume sent : %s", email),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
