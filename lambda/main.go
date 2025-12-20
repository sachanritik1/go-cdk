package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandlerRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("Username is required")

	}
	return fmt.Sprintf("Hello, %s!", event.Username), nil
}

func main() {
	lambda.Start(HandlerRequest)
}
