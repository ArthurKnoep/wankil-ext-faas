package handler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/ArthurKnoep/wankil-ext-token-faas/function/config"
	"github.com/ArthurKnoep/wankil-ext-token-faas/function/token"
	"github.com/ArthurKnoep/wankil-ext-token-faas/function/twitch"
)

func HandleRequest(config *config.Config, ddb *dynamodb.DynamoDB) func(ctx context.Context, event events.SQSEvent) (string, error) {
	return func(ctx context.Context, event events.SQSEvent) (string, error) {
		t, err := token.GetToken(config, ddb)
		if err != nil {
			return "", err
		}
		streams, err := twitch.GetStreams(config.StreamerIds, config, t.Token)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(streams)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}
