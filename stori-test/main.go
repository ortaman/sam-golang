package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ortaman/stori-test/adapters"
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

	// Validate that email exis in the body request
	bodyMap := map[string]string{}

	if err := json.Unmarshal([]byte(request.Body), &bodyMap); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	email, emailExist := bodyMap["email"]

	if !emailExist {
		return events.APIGatewayProxyResponse{
			Body:       "Email field is required",
			StatusCode: 400,
		}, nil
	}

	// Read csv file
	csvLinesPointer, err := adapters.ReadAllCSV(csvDir)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	// Validate csv file and load data in a slice
	transactionData, err := adapters.ValidateTransactionsFromCSV(csvLinesPointer)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	db := infra.NewMyPSQLConnection()
	PSQConfig := infra.NewPSQLConfig()

	dbRepository := repository.NewSQLRepository(db)
	emailRepository := repository.NewEmailRepository(PSQConfig)

	tnxsUsercase := usecase.NewTnxsUsecase(dbRepository, emailRepository)

	tnxsUsercase.SaveTransactions(email, &transactionData)
	tnxsUsercase.SendEmail(&transactionData, []string{email}, templateDir)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Transactions Resume sent : %s", email),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
