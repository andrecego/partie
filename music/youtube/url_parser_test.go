package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchYoutube(t *testing.T) {
	assert.True(t, MatchURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ"))
	assert.True(t, MatchURL("youtube.com/watch?v=dQw4w9WgXcQ"))
	assert.True(t, MatchURL("https://youtu.be/_dWp3ZbP_DA"))
	assert.True(t, MatchURL("youtu.be/_dWp3ZbP_DA"))
	assert.False(t, MatchURL("youtube"))
	assert.False(t, MatchURL("youtu.be"))
}

func TestYoutubePlaylist(t *testing.T) {
	assert.True(t, MatchPlaylist("https://www.youtube.com/watch?v=fJ9rUzIMcZQ&list=RDEMbHaAxpOZhcVmmF6I3y0siA"))
	assert.True(t, MatchPlaylist("https://www.youtube.com/playlist?list=PLt-N5ZTwt4xmHElkkSSo-TWNG0yrBJr7s"))
	assert.True(t, MatchPlaylist("https://youtu.be/azdwsXLmrHE?list=RDEMbHaAxpOZhcVmmF6I3y0siA&t=56"))
	assert.True(t, MatchPlaylist("youtu.be/azdwsXLmrHE?list=RDEMbHaAxpOZhcVmmF6I3y0siA&t=56"))
	assert.False(t, MatchPlaylist("https://www.youtube.com/watch?v=fJ9rUzIMcZQ"))
	assert.False(t, MatchPlaylist("https://www.youtube.com/list?v=listRQQ"))
}
