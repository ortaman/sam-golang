package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ortaman/stori-test/usecase"

	_ "github.com/lib/pq"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	csvDir := "./txns.csv"
	transactions, err := usecase.LoadTransactions(csvDir)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	transactionsResume, _ := usecase.GetTransactionsResume(transactions)
	templateDir := "./txns_template.html"

	if err := usecase.SendEmail([]string{email}, templateDir, &transactionsResume); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
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
