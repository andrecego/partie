package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"partie-bot/cache"
	"partie-bot/config"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	spotifyTokenKey = "spotify_access_token"
)

func SpotifyGet(endpoint string, queryParams map[string]string) (*http.Response, error) {
	token, err := fetchToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, _ := http.NewRequest("GET", endpoint, nil)
	request.Header.Set("Authorization", "Bearer "+token)

	q := request.URL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}

	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)
	return response, err
}

type authenticateResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func fetchToken() (string, error) {
	redisClient := cache.New().Client
	token, err := redisClient.Get(context.TODO(), spotifyTokenKey).Result()
	if err == nil {
		return token, nil
	}

	if err != redis.Nil {
		fmt.Println("Error fetching spotify token: ", err)
		return "", err
	}

	url := "https://accounts.spotify.com/api/token"

	payload := strings.NewReader("grant_type=client_credentials")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "partie/0.0.0")
	req.Header.Add("Authorization", "Basic "+config.SpotifyBearer)

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	var response authenticateResponse
	json.NewDecoder(res.Body).Decode(&response)

	cmd := redisClient.Set(context.TODO(), spotifyTokenKey, response.AccessToken, time.Duration(response.ExpiresIn)*time.Second)
	if cmd.Err() != nil {
		fmt.Println("Error saving spotify token: ", cmd.Err())
		return "", cmd.Err()
	}

	return response.AccessToken, nil
}
