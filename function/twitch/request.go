package twitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ArthurKnoep/wankil-ext-token-faas/function/config"
)

type (
	TokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   uint   `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	Requester struct {
		client http.Client
		conf *config.Config
	}
)

func NewRequester(conf *config.Config) *Requester {
	return &Requester{
		client: http.Client{
			Timeout: 30 * time.Second,
		},
		conf: conf,
	}
}

func (r *Requester) GetToken(config *config.Config) (*TokenResp, error) {
	u, err := url.Parse("https://id.twitch.tv/oauth2/token")
	if err != nil {
		return nil, err
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
		return nil, errors.New("could not request twitch api")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New("invalid response code from twitch api")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("could not read response from twitch api")
	}
	var t TokenResp
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Requester) GetStreams(token string) ([]Stream, error) {
	u, err := url.Parse("https://api.twitch.tv/helix/streams")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for _, streamerId := range r.conf.StreamerIds {
		q.Add("user_id", streamerId)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-ID", r.conf.ClientID)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New("invalid response code from twitch api")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var streamsReq StreamsRequest
	if err := json.Unmarshal(data, &streamsReq); err != nil {
		return nil, err
	}
	return streamsReq.Data, nil
}

func (r *Requester) GetGames(gameIds []string, token string) ([]Game, error) {
	u, err := url.Parse("https://api.twitch.tv/helix/games")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	for _, gameId := range gameIds {
		q.Add("id", gameId)
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Client-ID", r.conf.ClientID)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, errors.New("invalid response code from twitch api")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var gamesReq GamesRequest
	if err := json.Unmarshal(data, &gamesReq); err != nil {
		return nil, err
	}
	return gamesReq.Data, nil
}
