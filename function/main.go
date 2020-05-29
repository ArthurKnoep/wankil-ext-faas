package main

import (
	"fmt"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/ArthurKnoep/wankil-ext-token-faas/function/config"
	"github.com/ArthurKnoep/wankil-ext-token-faas/function/handler"
)

func main() {
	conf, err := config.Parse()
	if err != nil {
		panic(err)
	}
	var ddb *dynamodb.DynamoDB
	if sess, err := session.NewSession(&aws.Config{
		Region: &conf.Region,
	}); err != nil {
		fmt.Printf("Failed to connect to DynamoDB: %s\n", err.Error())
		panic(err)
	} else {
		ddb = dynamodb.New(sess)
	}
	runtime.Start(handler.HandleRequest(conf, ddb))
}
