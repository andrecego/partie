package youtube

import (
	"fmt"
	"partie-bot/system"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type YoutubeResult struct {
	Entries []Youtube `json:"entries"`
}

type Youtube struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	VideoURL  string `json:"webpage_url"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	StartTime int
	AddedBy   AddedBy
}

type AddedBy struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
}

func (y *Youtube) GetID() string {
	return y.ID
}

func (y *Youtube) GetTitle() string {
	titleNomalizedRegex := regexp.MustCompile(`(?P<Title>.*[-:].*)([\(\[\|].+)$`)
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
	return fmt.Sprintf("<@%s>", y.AddedBy.User.ID)
}

func (y *Youtube) GetAuthorID() string {
	return y.AddedBy.User.ID
}

func (y *Youtube) GetChannelID() string {
	return y.AddedBy.Channel.ID
}

func (y *Youtube) GetGuildID() string {
	return y.AddedBy.Guild.ID
}

func (y *Youtube) GetURL() string {
	if MatchURL(y.URL) {
		getUrlArgs := append(youtubeFindUrlArgs, y.URL)
		err, stdout, stderr := system.ShellOut(strings.Join(getUrlArgs, " "))
		if err != nil {
			fmt.Println("Error getting url from youtube: ", err)
			fmt.Println(stderr)
			return ""
		}

		return strings.TrimSpace(stdout)
	}

	return y.URL
}

func (y *Youtube) GetThumbnail() string {
	return y.Thumbnail
}

func (y *Youtube) GetStartTime() int {
	return y.StartTime
}

func (y *Youtube) GetVideoURL() string {
	return y.VideoURL
}

var youtubeDefaultArgs = []string{
	"/usr/local/bin/yt-dlp",
	"--no-playlist",
	"--dump-single-json",
	"-x",
}

var youtubeFindUrlArgs = []string{
	"/usr/local/bin/yt-dlp",
	"--no-playlist",
	"-x",
	"--get-url",
}

var youtubePlaylistArgs = []string{
	"/usr/local/bin/yt-dlp",
	"--flat-playlist",
	"--dump-single-json",
	"-x",
	"--playlist-end 20",
}
