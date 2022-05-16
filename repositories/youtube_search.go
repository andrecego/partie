package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"partie-bot/cache"
	"partie-bot/requests"
	"strconv"
	"strings"
)

type YoutubeSearchResult struct {
	Title       string `json:"title"`
	URL         string `json:"webpage_url"`
	DurationRaw string `json:"duration_raw"`
}

func (r YoutubeSearchResult) Duration() int {
	duration := 0

	time := strings.Split(r.DurationRaw, ":")
	timeLen := len(time)
	if timeLen > 3 {
		return 0
	}

	for i, t := range time {
		timeInt, err := strconv.Atoi(t)
		if err != nil {
			return 0
		}

		power := timeLen - i - 1

		duration += timeInt * int(math.Pow(60, float64(power)))
	}

	return duration
}

func YoutubeSearch(query string) ([]YoutubeSearchResult, error) {
	resp, err := requests.YoutubeSearch(query)
	if err != nil {
		return nil, fmt.Errorf("Error searching youtube: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error searching youtube: %s", resp.Status)
	}

	var results []YoutubeSearchResult
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling youtube search response: %s", err)
	}

	return results, nil
}

func Save(guildID, userID string, searchResults []YoutubeSearchResult) error {
	redisClient := cache.New().Client

	key := fmt.Sprintf("guilds:%s:users:%s:search_results", guildID, userID)

	searchBytes, err := json.Marshal(searchResults)
	if err != nil {
		return fmt.Errorf("Error marshalling search results: %v", err)
	}

	err = redisClient.Set(context.TODO(), key, searchBytes, 0).Err()
	if err != nil {
		return fmt.Errorf("Error saving youtube search results: %v", err)
	}

	return nil
}
