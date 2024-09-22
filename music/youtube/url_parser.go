package youtube

import "regexp"

func MatchBaseURL(query string) bool {
	return regexp.
		MustCompile(`^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+$`).
		Match([]byte(query))
}

func CaptureVideoId(query string) string {
	videoId := regexp.
		MustCompile(`(?:v=|youtu\.be\/)([^&?]+)`).
		FindSubmatch([]byte(query))
	if len(videoId) < 2 {
		return ""
	}

	return string(videoId[1])
}

func MatchPlaylist(query string) bool {
	return regexp.
		MustCompile(`^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.*list=.+$`).
		Match([]byte(query))
}
