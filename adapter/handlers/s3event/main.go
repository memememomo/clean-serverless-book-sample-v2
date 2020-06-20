package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(event events.S3Event) error {
	fmt.Printf("%+v", event)
	return nil
}

func main() {
	lambda.Start(handler)
}
