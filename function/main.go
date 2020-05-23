package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return "", err
	}
	lc, _ := lambdacontext.FromContext(ctx)
	resp := string(data)
	resp += lc.AwsRequestID
	return resp, nil
}

func main() {
	runtime.Start(handleRequest)
}
