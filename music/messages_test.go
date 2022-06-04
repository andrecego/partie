package music

import (
	"partie-bot/music/youtube"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

type testGetMessageCase struct {
	description          string
	currentDJ            *DJ
	expectedQueueMessage string
	expectedNPStatus     *discordgo.MessageEmbed
}

var testGetMessageCases = []testGetMessageCase{
	{
		description: "When there is no song playing",
		currentDJ: &DJ{
			CurrentSong: nil,
			Queue:       nil,
		},
		expectedQueueMessage: "__**Queue:**__\nNo songs in queue. Go add some! 🎵",
		expectedNPStatus: &discordgo.MessageEmbed{
			Title: "Nothing playing.",
		},
	},
	{
		description: "When there is some songs in the queue",
		currentDJ: &DJ{
			CurrentSong: nil,
			Queue: []Song{
				&youtube.Youtube{Title: "Song 1", Duration: 1, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "1"}}},
				&youtube.Youtube{Title: "Song 2", Duration: 2, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "2"}}},
			},
		},
		expectedQueueMessage: "" +
			"__**Queue:**__\n" +
			"02. Song 2 [00:02] - <@2>\n" +
			"01. Song 1 [00:01] - <@1>\n",
		expectedNPStatus: &discordgo.MessageEmbed{
			Title: "Nothing playing.",
		},
	},
	{
		description: "When there is many songs in the queue",
		currentDJ: &DJ{
			CurrentSong: nil,
			Queue: []Song{
				&youtube.Youtube{Title: "Song 1", Duration: 1, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "1"}}},
				&youtube.Youtube{Title: "Song 2", Duration: 2, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "2"}}},
				&youtube.Youtube{Title: "Song 3", Duration: 3, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "3"}}},
				&youtube.Youtube{Title: "Song 4", Duration: 4, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "4"}}},
				&youtube.Youtube{Title: "Song 5", Duration: 5, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "5"}}},
				&youtube.Youtube{Title: "Song 6", Duration: 6, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "6"}}},
				&youtube.Youtube{Title: "Song 7", Duration: 7, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "7"}}},
				&youtube.Youtube{Title: "Song 8", Duration: 8, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "8"}}},
				&youtube.Youtube{Title: "Song 9", Duration: 9, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "9"}}},
				&youtube.Youtube{Title: "Song 10", Duration: 10, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "10"}}},
				&youtube.Youtube{Title: "Song 11", Duration: 11, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "11"}}},
				&youtube.Youtube{Title: "Song 12", Duration: 12, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "12"}}},
				&youtube.Youtube{Title: "Song 13", Duration: 13, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "13"}}},
				&youtube.Youtube{Title: "Song 14", Duration: 14, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "14"}}},
				&youtube.Youtube{Title: "Song 15", Duration: 15, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "15"}}},
				&youtube.Youtube{Title: "Song 16", Duration: 16, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "16"}}},
				&youtube.Youtube{Title: "Song 17", Duration: 17, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "17"}}},
				&youtube.Youtube{Title: "Song 18", Duration: 18, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "18"}}},
				&youtube.Youtube{Title: "Song 19", Duration: 19, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "19"}}},
				&youtube.Youtube{Title: "Song 20", Duration: 20, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "20"}}},
				&youtube.Youtube{Title: "Song 21", Duration: 21, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "21"}}},
				&youtube.Youtube{Title: "Song 22", Duration: 22, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "22"}}},
				&youtube.Youtube{Title: "Song 23", Duration: 23, AddedBy: youtube.AddedBy{User: &discordgo.User{ID: "23"}}},
			},
		},
		expectedQueueMessage: "" +
			"__**Queue:**__\n" +
			"And 3 more songs\n" +
			"20. Song 20 [00:20] - <@20>\n" +
			"19. Song 19 [00:19] - <@19>\n" +
			"18. Song 18 [00:18] - <@18>\n" +
			"17. Song 17 [00:17] - <@17>\n" +
			"16. Song 16 [00:16] - <@16>\n" +
			"15. Song 15 [00:15] - <@15>\n" +
			"14. Song 14 [00:14] - <@14>\n" +
			"13. Song 13 [00:13] - <@13>\n" +
			"12. Song 12 [00:12] - <@12>\n" +
			"11. Song 11 [00:11] - <@11>\n" +
			"10. Song 10 [00:10] - <@10>\n" +
			"09. Song 9 [00:09] - <@9>\n" +
			"08. Song 8 [00:08] - <@8>\n" +
			"07. Song 7 [00:07] - <@7>\n" +
			"06. Song 6 [00:06] - <@6>\n" +
			"05. Song 5 [00:05] - <@5>\n" +
			"04. Song 4 [00:04] - <@4>\n" +
			"03. Song 3 [00:03] - <@3>\n" +
			"02. Song 2 [00:02] - <@2>\n" +
			"01. Song 1 [00:01] - <@1>\n",
		expectedNPStatus: &discordgo.MessageEmbed{
			Title: "Nothing playing.",
		},
	},
}

func TestGetMessage(t *testing.T) {
	for _, testCase := range testGetMessageCases {
		t.Run(testCase.description, func(t *testing.T) {
			currentDJ = testCase.currentDJ

			queueMessage, npStatus := getMessage()
			assert.Equal(t, testCase.expectedQueueMessage, queueMessage)
			assert.Equal(t, testCase.expectedNPStatus.Title, npStatus.Title)
		})
	}
}
