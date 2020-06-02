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

func appendIfNotExist(slice []string, toAppend string) []string {
	for _, elem := range slice {
		if elem == toAppend {
			return slice
		}
	}
	return append(slice, toAppend)
}

func findGame(games []twitch.Game, gameId string) twitch.Game {
	for _, game := range games {
		if game.Id == gameId {
			return game
		}
	}
	return twitch.Game{
		Id:        "0",
		Name:      "Unknown game",
		BoxArtUrl: "",
	}
}

func createStreamsObject(streams []twitch.Stream, games []twitch.Game) Streams {
	var resp Streams
	for _, stream := range streams {
		game := findGame(games, stream.GameId)
		resp.Streams = append(resp.Streams, Stream{
			Stream: twitch.Stream{
				Id:           stream.Id,
				UserId:       stream.UserId,
				UserName:     stream.UserName,
				GameId:       stream.GameId,
				Type:         stream.Type,
				Title:        stream.Title,
				ViewerCount:  stream.ViewerCount,
				StartedAt:    stream.StartedAt,
				Language:     stream.Language,
				ThumbnailUrl: stream.ThumbnailUrl,
				TagIds:       stream.TagIds,
			},
			GameName:      game.Name,
			GameBoxArtUrl: game.BoxArtUrl,
		})
	}
	return resp
}

func HandleRequest(config *config.Config, ddb *dynamodb.DynamoDB) func(ctx context.Context, event events.SQSEvent) (string, error) {
	return func(ctx context.Context, event events.SQSEvent) (string, error) {
		t, err := token.GetToken(config, ddb)
		if err != nil {
			return "", err
		}
		r := twitch.NewRequester(config)
		streams, err := r.GetStreams(t.Token)
		if err != nil {
			return "", err
		}
		var gameIdNeeded []string
		for _, stream := range streams {
			gameIdNeeded = appendIfNotExist(gameIdNeeded, stream.GameId)
		}
		var games []twitch.Game
		if len(gameIdNeeded) > 0 {
			games, err = r.GetGames(gameIdNeeded, t.Token)
			if err != nil {
				return "", err
			}
		}
		response := createStreamsObject(streams, games)
		bytes, err := json.Marshal(response)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}
