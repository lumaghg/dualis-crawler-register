package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	tableName = "DUALIS_CREDENTIALS"
)

type MyEvent struct {
	Email             string `json:"email"`
	NotificationEmail string `json:"notificationEmail"`
	Password          string `json:"password"`
}

type User struct {
	Email             string
	NotificationEmail string
	Password          string
}

func HandleRequest(ctx context.Context, request MyEvent) (int, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(request.Email),
			},
		},
	})
	if err != nil {
		return 500, err
	}

	if result != nil {
		return 400, fmt.Errorf("User with email %v already registered", request.Email)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(request.Email),
			},
			"NotificationEmail": {
				S: aws.String(request.NotificationEmail),
			},
			"Password": {
				S: aws.String(request.Password),
			},
		},
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		return 500, err
	}
	return 201, nil
}

func main() {
	lambda.Start(HandleRequest)
}
