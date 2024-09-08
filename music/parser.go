package music

import (
	"fmt"
	"partie-bot/interfaces"
	"partie-bot/music/spotify"
	"partie-bot/music/youtube"
)

func ParseQuery(query string) []interfaces.Finder {
	finders := parseYoutubeQuery(query)
	if finders != nil {
		return finders
	}

	finders = parseSpotifyQuery(query)
	if finders != nil {
		return finders
	}

	return []interfaces.Finder{youtube.TextParse(query)}
}

func parseYoutubeQuery(query string) []interfaces.Finder {
	if !youtube.MatchURL(query) {
		return nil
	}

	if youtube.MatchPlaylist(query) {
		return []interfaces.Finder{youtube.PlaylistParse(query)}
	}

	return []interfaces.Finder{youtube.TextParse(query)}
}

func parseSpotifyQuery(query string) []interfaces.Finder {
	if !spotify.MatchURL(query) {
		return nil
	}

	if spotify.MatchTrack(query) {
		return []interfaces.Finder{youtube.TextParse(query)}
	}

	if spotify.MatchPlaylist(query) {
		playlistID := spotify.PlaylistIdCapture(query)
		if playlistID == "" {
			fmt.Println("Error parsing spotify playlist ID")
			return nil
		}

		songTitles := spotify.PlaylistURLToTexts(playlistID)
		if len(songTitles) == 0 {
			fmt.Println("No tracks found in the playlist")
			return nil
		}
		return youtube.TextParseSlice(songTitles)
	}

	fmt.Println("No match found for spotify query: ", query)

	return nil
}
