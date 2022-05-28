package youtube

import (
	"fmt"
	"partie-bot/interfaces"
	"partie-bot/system"
	"strings"
)

type textFinder struct {
	Query string
}

func (f textFinder) GetQuery() string {
	return f.Query
}

func (f textFinder) Download() (string, error) {
	if f.Query == "" {
		return "", ErrorNoQuery
	}

	err, stdout, stderr := system.ShellOut(textDownloadCommand(f.Query + " lyrics"))
	if err != nil {
		fmt.Println("Error downloading of youtube: ", err)
		fmt.Println(stderr)
		return "", err
	}

	return strings.TrimSpace(stdout), nil
}

func TextParse(query string) interfaces.Finder {
	return textFinder{Query: query}
}

func textDownloadCommand(query string) string {
	dlArgs := append(youtubeDefaultArgs, fmt.Sprintf(`"ytsearch:%s"`, query))

	return strings.Join(dlArgs, " ")
}
