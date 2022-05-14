package music

import (
	"partie-bot/interfaces"
	"partie-bot/music/youtube"
)

func ParseQuery(query string) interfaces.Finder {
	if youtube.MatchURL(query) {
		if youtube.MatchPlaylist(query) {
			return youtube.PlaylistParse(query)
		}

		return youtube.URLParse(query)
	}

	return youtube.TextParse(query)
}
