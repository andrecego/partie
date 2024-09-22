package youtube

import (
	"errors"
	"fmt"
	"strings"
)

var ErrorOnlyMusicPremium = errors.New("%s is only available to Music Premium members")

func (y urlFinder) handleYoutubeDLPError(err error, stdout, stderr string) (string, error) {
	if strings.Contains(stderr, "This video is only available to Music Premium members") {
		return "", fmt.Errorf(ErrorOnlyMusicPremium.Error(), y.Query)
	}

	fmt.Println("Error downloading youtube from URL: ", err)
	fmt.Println(stderr)

	return "", err
}
