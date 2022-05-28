package spotify

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Mock the external request
func TestURLToText(t *testing.T) {
	type testCase struct {
		description  string
		url          string
		expectedText string
	}

	testCases := []testCase{
		{
			description:  "Spotify track URL",
			url:          "https://open.spotify.com/track/2LEF1A8DOZ9wRYikWgVlZ8",
			expectedText: "Good Feeling - song by Flo Rida | Spotify",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedText, URLToText(testCase.url), testCase.description)
	}
}
