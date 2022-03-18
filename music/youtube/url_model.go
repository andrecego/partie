package youtube

import (
	"fmt"
	"partie-bot/interfaces"
	"partie-bot/system"
	"strings"
)

type urlFinder struct {
	Query string
}

func (y urlFinder) GetQuery() string {
	return y.Query
}

func (y urlFinder) Download() (string, error) {
	err, stdout, stderr := system.ShellOut(urlDownloadCommand(y.Query))
	if err != nil {
		fmt.Println("Error downloading youtube from URL: ", err)
		fmt.Println(stderr)
		return "", err
	}

	return strings.TrimSpace(stdout), nil
}

func URLParse(query string) interfaces.Finder {
	return urlFinder{Query: query}
}

func urlDownloadCommand(url string) string {
	dlArgs := append(youtubeDefaultArgs, url)

	return strings.Join(dlArgs, " ")
}
