package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/lib/pq"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println(request)

	return events.APIGatewayProxyResponse{
		Body:       "Hello Stori!",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
