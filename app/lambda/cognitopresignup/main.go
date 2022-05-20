package main

import (
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler is the lambda handler invoked by the `lambda.Start` function call
func handler(event events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	v, ok := event.Request.UserAttributes["custom:isActive"]
	if !ok {
		return event, nil
	}
	autoConfirmUser, err := strconv.ParseBool(v)
	if err != nil {
		return event, nil
	}
	event.Response.AutoConfirmUser = autoConfirmUser
	return event, nil
}

func main() {
	lambda.Start(handler)
}
