package spotify

import "regexp"

func MatchURL(query string) bool {
	return regexp.MustCompile(`^(https?:\/\/)?(www\.)?(open\.spotify\.com)\/.+$`).Match([]byte(query))
}

func MatchTrack(query string) bool {
	return regexp.MustCompile(`^(https?:\/\/)?(www\.)?(open\.spotify\.com)\/track\/.+$`).Match([]byte(query))
}

func MatchPlaylist(query string) bool {
	return regexp.MustCompile(`^(https?:\/\/)?(www\.)?(open\.spotify\.com)\/playlist\/.+$`).Match([]byte(query))
}

func PlaylistIdCapture(query string) string {
	playlistId := regexp.MustCompile(`^(https?:\/\/)?(www\.)?(open\.spotify\.com)\/playlist\/([^?]+)`).FindStringSubmatch(query)
	if len(playlistId) == 0 {
		return ""
	}

	return playlistId[4]
}
