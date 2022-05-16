package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDurationCase struct {
	description string
	durationRaw string
	expected    int
}

var testDurationCases = []testDurationCase{
	{
		description: "with only seconds",
		durationRaw: "01",
		expected:    01,
	},
	{
		description: "with minutes and seconds",
		durationRaw: "26:59",
		expected:    26*60 + 59,
	},
	{
		description: "with hours minutes and seconds",
		durationRaw: "01:02:03",
		expected:    01*60*60 + 02*60 + 03,
	},
	{
		description: "with days, hours minutes and seconds",
		durationRaw: "04:01:02:03",
		expected:    0,
	},
	{
		description: "when duration is invalid",
		durationRaw: "1h20m",
		expected:    0,
	},
	{
		description: "when duration is zero",
		durationRaw: "0",
		expected:    0,
	},
}

func TestDuration(t *testing.T) {

	for _, testCase := range testDurationCases {
		youtubeSearchResult := YoutubeSearchResult{DurationRaw: testCase.durationRaw}
		result := youtubeSearchResult.Duration()

		assert.Equal(t, testCase.expected, result, testCase.description)
	}
}
