package music

import "fmt"

func nowPlayingMessage() {
	channelID := "943655307626823771"
	song := currentDJ.CurrentSong

	message := fmt.Sprintf("Now playing: %s\nAdded by: %s", song.GetTitle(), song.GetAddedBy())
	currentDJ.Discord.Session.ChannelMessageSend(channelID, message)
}

func queueMessage(channelID string) {
	message := ""

	np := currentDJ.CurrentSong
	if np != nil {
		message += fmt.Sprintf("Now playing: %s - %s\n", np.GetTitle(), np.GetAddedBy())
	}

	if len(currentDJ.Queue) > 0 {
		message += "Queue:\n"
		for i, song := range currentDJ.Queue {
			message += fmt.Sprintf("%d. %s - %s\n", i+1, song.GetTitle(), song.GetAddedBy())
		}
	}

	if message == "" {
		message = "No songs in queue. Go add some! 🎵"
	}

	currentDJ.Discord.Session.ChannelMessageSend(channelID, message)
}

func addedToQueueMessage(song Song, channelID string) {
	message := fmt.Sprintf("Added %s to queue.\n", song.GetTitle())
	currentDJ.Discord.Session.ChannelMessageSend(channelID, message)
}
