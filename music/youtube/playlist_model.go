package youtube

import (
	"fmt"
	"partie-bot/interfaces"
	"partie-bot/system"
	"strings"
)

type playlistFinder struct {
	Query string
}

func (y playlistFinder) GetQuery() string {
	return y.Query
}

func (y playlistFinder) Download() (string, error) {
	err, stdout, stderr := system.ShellOut(playlistDownloadCommand(y.Query))
	if err != nil {
		fmt.Println("Error downloading youtube playlist from URL: ", err)
		fmt.Println(stderr)
		return "", err
	}

	return strings.TrimSpace(stdout), nil
}

func PlaylistParse(query string) interfaces.Finder {
	return playlistFinder{Query: query}
}

func playlistDownloadCommand(url string) string {
	dlArgs := append(youtubePlaylistArgs, url)

	return strings.Join(dlArgs, " ")
}
