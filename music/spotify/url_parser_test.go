package spotify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchURL(t *testing.T) {
	assert.True(t, MatchURL("https://open.spotify.com/track/2LEF1A8DOZ9wRYikWgVlZ8?si=1817bbe6fa7f4c7b"))
	assert.True(t, MatchURL("https://open.spotify.com/playlist/1rqGgM1fIIiS3nDFUmdk2g"))
	assert.True(t, MatchURL("open.spotify.com/track/2LEF1A8DOZ9wRYikWgVlZ8"))
	assert.True(t, MatchURL("https://open.spotify.com/something"))
	assert.False(t, MatchURL("open.spotify.com"))
	assert.False(t, MatchURL("https://youtu.be/azdwsXLmrHE"))
}

func TestMatchTrack(t *testing.T) {
	assert.True(t, MatchTrack("https://open.spotify.com/track/2LEF1A8DOZ9wRYikWgVlZ8?si=1817bbe6fa7f4c7b"))
	assert.True(t, MatchTrack("http://open.spotify.com/track/2LEF1A8DOZ9wRYikWgVlZ8"))
	assert.True(t, MatchTrack("open.spotify.com/track/2LEF1A8DOZ9wRYikWgVlZ8"))
	assert.False(t, MatchTrack("https://open.spotify.com/playlist/1rqGgM1fIIiS3nDFUmdk2g"))
	assert.False(t, MatchTrack("https://open.spotify.com/something"))
	assert.False(t, MatchTrack("open.spotify.com"))
}
