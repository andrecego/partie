package youtube

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchYoutubeBaseURL(t *testing.T) {
	assert.True(t, MatchBaseURL("https://www.youtube.com/watch?v=dQw4w9WgXcQ"))
	assert.True(t, MatchBaseURL("youtube.com/watch?v=dQw4w9WgXcQ"))
	assert.True(t, MatchBaseURL("https://youtu.be/_dWp3ZbP_DA"))
	assert.True(t, MatchBaseURL("youtu.be/_dWp3ZbP_DA"))
	assert.False(t, MatchBaseURL("youtube"))
	assert.False(t, MatchBaseURL("youtu.be"))
}

func TestCaptureVideoId(t *testing.T) {
	assert.Equal(t, "dQw4w9WgXcQ", CaptureVideoId("https://www.youtube.com/watch?v=dQw4w9WgXcQ"))
	assert.Equal(t, "dQw4w9WgXcQ", CaptureVideoId("youtube.com/watch?v=dQw4w9WgXcQ"))
	assert.Equal(t, "_dWp3ZbP_DA", CaptureVideoId("https://youtu.be/_dWp3ZbP_DA"))
	assert.Equal(t, "_dWp3ZbP_DA", CaptureVideoId("youtu.be/_dWp3ZbP_DA"))
}

func TestYoutubePlaylist(t *testing.T) {
	assert.True(t, MatchPlaylist("https://www.youtube.com/watch?v=fJ9rUzIMcZQ&list=RDEMbHaAxpOZhcVmmF6I3y0siA"))
	assert.True(t, MatchPlaylist("https://www.youtube.com/playlist?list=PLt-N5ZTwt4xmHElkkSSo-TWNG0yrBJr7s"))
	assert.True(t, MatchPlaylist("https://youtu.be/azdwsXLmrHE?list=RDEMbHaAxpOZhcVmmF6I3y0siA&t=56"))
	assert.True(t, MatchPlaylist("youtu.be/azdwsXLmrHE?list=RDEMbHaAxpOZhcVmmF6I3y0siA&t=56"))
	assert.False(t, MatchPlaylist("https://www.youtube.com/watch?v=fJ9rUzIMcZQ"))
	assert.False(t, MatchPlaylist("https://www.youtube.com/list?v=listRQQ"))
}
