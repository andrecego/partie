package music

import (
	"partie-bot/interfaces"
	"partie-bot/music/spotify"
	"partie-bot/music/youtube"
)

func ParseQuery(query string) interfaces.Finder {
	if youtube.MatchURL(query) {
		if youtube.MatchPlaylist(query) {
			return youtube.PlaylistParse(query)
		}

		return youtube.URLParse(query)
	}

	if spotify.MatchURL(query) {
		if spotify.MatchTrack(query) {
			query = spotify.URLToText(query)
		} else {
			query = ""
		}
	}

	return youtube.TextParse(query)
}
