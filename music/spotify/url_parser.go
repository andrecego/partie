package spotify

import "regexp"

func MatchURL(query string) bool {
	return regexp.MustCompile(`^(https?:\/\/)?(www\.)?(open\.spotify\.com)\/.+$`).Match([]byte(query))
}

func MatchTrack(query string) bool {
	return regexp.MustCompile(`^(https?:\/\/)?(www\.)?(open\.spotify\.com)\/track\/.+$`).Match([]byte(query))
}
