package music

import (
	"partie-bot/interfaces"
	"partie-bot/music/youtube"
)

func ParseQuery(query string) interfaces.Finder {
	if youtube.MatchURL(query) {
		return youtube.URLParse(query)
	}

	return youtube.TextParse(query)
}
