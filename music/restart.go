package music

import (
	"github.com/bwmarrin/discordgo"
)

func Restart(session *discordgo.Session, guildID, userID string) {
	if currentDJ == nil {
		return
	}

	oldQueue := append([]Song{currentDJ.CurrentSong}, currentDJ.Queue...)

	currentDJ = nil
	New(session)

	currentDJ.Queue = oldQueue
	currentDJ.Paused = false
	currentDJ.NeedsToSkip = true

	// need to call stream if is not playing
	// Stream(session, &discordgo.MessageCreate{Message: &discordgo.Message{GuildID: guildID, Author: &discordgo.User{ID: userID}}})
}
