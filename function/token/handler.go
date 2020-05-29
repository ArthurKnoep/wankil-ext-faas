package token

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/ArthurKnoep/wankil-ext-token-faas/function/config"
	"github.com/ArthurKnoep/wankil-ext-token-faas/function/twitch"
)

const tableName = "wankil-ext-token"

type Token struct {
	Id        int64     `json:"id"`
	Token     string    `json:"token"`
	Expiry    time.Time `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
}

func getTokenFromDB(ddb *dynamodb.DynamoDB) (*Token, error) {
	rst, err := ddb.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	if rst.Count == aws.Int64(0) {
		return nil, nil
	}
	for _, v := range rst.Items {
		token := Token{}
		if err := dynamodbattribute.UnmarshalMap(v, &token); err != nil {
			return nil, err
		}
		if token.Expiry.After(time.Now()) {
			return &token, nil
		}
	}
	return nil, nil
}

func getTokenFromTwitch(conf *config.Config, ddb *dynamodb.DynamoDB) (*Token, error) {
	resp, err := twitch.GetToken(conf)
	if err != nil {
		return nil, err
	}
	token := Token{
		Id:        1,
		Token:     resp.AccessToken,
		Expiry:    time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second),
		CreatedAt: time.Now(),
	}
	item, err := dynamodbattribute.MarshalMap(token)
	if err != nil {
		return nil, err
	}
	if _, err := ddb.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}); err != nil {
		return nil, err
	}
	return &token, nil
}

func GetToken(conf *config.Config, ddb *dynamodb.DynamoDB) (*Token, error) {
	token, err := getTokenFromDB(ddb)
	if err != nil {
		return nil, err
	}
	if token != nil {
		return token, nil
	}
	token, err = getTokenFromTwitch(conf, ddb)
	if err != nil {
		return nil, err
	}
	return token, nil
}
