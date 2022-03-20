package youtube

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Youtube struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	URL      string `json:"url"`
	AddedBy  *discordgo.User
}

func (y *Youtube) GetID() string {
	return y.ID
}

func (y *Youtube) GetTitle() string {
	return y.Title
}

func (y *Youtube) GetDuration() time.Duration {
	return time.Duration(y.Duration) * time.Second
}

func (y *Youtube) GetAddedBy() string {
	return fmt.Sprintf("%s#%s", y.AddedBy.Username, y.AddedBy.Discriminator)
}

func (y *Youtube) GetURL() string {
	return y.URL
}

var youtubeDefaultArgs = []string{
	"/usr/local/bin/yt-dlp",
	"--no-playlist",
	"--dump-json",
	"-x",
}
