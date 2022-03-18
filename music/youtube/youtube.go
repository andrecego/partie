package youtube

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Youtube struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	URL      string `json:"webpage_url"`
	AddedBy  *discordgo.User
	Path     string
}

func (y *Youtube) GetID() string {
	return y.ID
}

func (y *Youtube) GetDuration() int {
	return y.Duration
}

func (y *Youtube) GetTitle() string {
	return y.Title
}

func (y *Youtube) GetAddedBy() string {
	return fmt.Sprintf("%s#%s", y.AddedBy.Username, y.AddedBy.Discriminator)
}

func (y *Youtube) GetPath() string {
	return y.Path
}

var youtubeDefaultArgs = []string{
	"/usr/local/bin/yt-dlp",
	"--no-playlist",
	"--write-info-json",
	"--print filename",
	"--no-simulate",
	"-q",
	"-x",
	`-o "tmp/youtube/%(id)s.opus"`,
}
