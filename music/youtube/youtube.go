package youtube

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Youtube struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	VideoURL  string `json:"webpage_url"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	AddedBy   *discordgo.User
}

func (y *Youtube) GetID() string {
	return y.ID
}

func (y *Youtube) GetTitle() string {
	titleNomalizedRegex := regexp.MustCompile(`(?P<Title>.*-.*)([\(\[\|].+)$`)
	matches := titleNomalizedRegex.FindStringSubmatch(y.Title)
	titleIndex := titleNomalizedRegex.SubexpIndex("Title")
	if titleIndex >= 0 && len(matches) > titleIndex {
		return matches[titleIndex]
	}

	return y.Title
}

func (y *Youtube) GetDuration() time.Duration {
	return time.Duration(y.Duration) * time.Second
}

func (y *Youtube) GetAddedBy() string {
	return fmt.Sprintf("<@%s>", y.AddedBy.ID)
}

func (y *Youtube) GetURL() string {
	return y.URL
}

func (y *Youtube) GetThumbnail() string {
	return y.Thumbnail
}

func (y *Youtube) GetVideoURL() string {
	return y.VideoURL
}

var youtubeDefaultArgs = []string{
	"/usr/local/bin/yt-dlp",
	"--no-playlist",
	"--dump-json",
	"-x",
}
