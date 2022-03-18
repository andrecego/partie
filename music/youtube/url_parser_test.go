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
