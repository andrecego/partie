package spotify

import (
	"encoding/json"
	"partie-bot/requests"

	"golang.org/x/net/html"
)

func URLToText(url string) string {
	response, err := requests.Get(url)
	if err != nil {
		return ""
	}

	tokenizer := html.NewTokenizer(response.Body)
	isTitle := false
	for {
		tt := tokenizer.Next()

		switch {
		case tt == html.ErrorToken:
			return ""

		case tt == html.StartTagToken:
			t := tokenizer.Token()
			isTitle = t.Data == "title"

		case tt == html.TextToken:
			t := tokenizer.Token()

			if isTitle {
				return t.Data
			}
		}
	}
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

type PlaylistResponse struct {
	Tracks struct {
		Items []struct {
			Track Track `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}

func PlaylistURLToTexts(url string) []string {
	queryParams := map[string]string{
		"fields": "tracks(items(track(name, artists(name))))",
	}

	response, err := requests.SpotifyGet(url, queryParams)
	if err != nil {
		return nil
	}

	var playlistResponse PlaylistResponse
	err = json.NewDecoder(response.Body).Decode(&playlistResponse)
	if err != nil {
		return nil
	}

	var texts []string
	for _, item := range playlistResponse.Tracks.Items {
		artists := ""
		for i, artist := range item.Track.Artists {
			if i > 0 {
				artists += ", "
			}
			artists += artist.Name
		}

		text := item.Track.Name + " - " + artists
		texts = append(texts, text)
	}

	return texts
}
