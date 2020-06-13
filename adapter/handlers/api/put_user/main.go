package main

import (
	"clean-serverless-book-sample-v2/adapter/controller"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return controller.PutUser(request), nil
}

func main() {
	lambda.Start(handler)
}
