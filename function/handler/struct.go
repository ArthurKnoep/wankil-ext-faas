package handler

import "github.com/ArthurKnoep/wankil-ext-token-faas/function/twitch"

type (
	Stream struct {
		twitch.Stream
		GameName      string `json:"game_name"`
		GameBoxArtUrl string `json:"game_box_art_url"`
	}

	Streams struct {
		Streams []Stream `json:"streams"`
	}
)
