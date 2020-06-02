package twitch

type (
	Stream struct {
		Id           string   `json:"id"`
		UserId       string   `json:"user_id"`
		UserName     string   `json:"user_name"`
		GameId       string   `json:"game_id"`
		Type         string   `json:"type"`
		Title        string   `json:"title"`
		ViewerCount  int64    `json:"viewer_count"`
		StartedAt    string   `json:"started_at"`
		Language     string   `json:"language"`
		ThumbnailUrl string   `json:"thumbnail_url"`
		TagIds       []string `json:"tag_ids"`
	}

	StreamsRequest struct {
		Data []Stream `json:"data"`
	}

	Game struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		BoxArtUrl string `json:"box_art_url"`
	}

	GamesRequest struct {
		Data []Game `json:"data"`
	}
)
