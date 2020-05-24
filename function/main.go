package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/caarlos0/env"
)

type Config struct {
	ClientID     string `env:"CLIENT_ID,required"`
	ClientSecret string `env:"CLIENT_SECRET,required"`
	Scope        string `env:"SCOPE"`
}

func handleRequest(config Config) func(ctx context.Context, event events.SQSEvent) (string, error) {
	return func(ctx context.Context, event events.SQSEvent) (string, error) {
		u, err := url.Parse("https://id.twitch.tv/oauth2/token")
		if err != nil {
			return "", err
		}
		q := u.Query()
		q.Add("client_id", config.ClientID)
		q.Add("client_secret", config.ClientSecret)
		q.Add("grant_type", "client_credentials")
		if len(config.Scope) > 0 {
			q.Add("scope", config.Scope)
		}
		u.RawQuery = q.Encode()
		resp, err := http.Post(u.String(), "", nil)
		if err != nil {
			return "", errors.New("could not request twitch api")
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 400 {
			return "", errors.New("invalid response code from twitch api")
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", errors.New("could not read response from twitch api")
		}
		return string(data), nil
	}
}

func main() {
	var config Config
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	runtime.Start(handleRequest(config))
}
