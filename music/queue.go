package music

import "fmt"

func ShowQueue(channelID string) {
	queueMessage(channelID)
}

func Skip() {
	currentDJ.NeedsToSkip = true
}

func Pause() {
	currentDJ.Paused = true
}

func Resume() {
	currentDJ.Paused = false
}

func skipSong() {
	currentDJ.NeedsToSkip = false
	if currentDJ.CurrentSong == nil {
		return
	}

	currentDJ.CurrentSong = nil
	Play()
}

func AddToQueue(song Song, channelID string) {
	currentDJ.Queue = append(currentDJ.Queue, song)
	addedToQueueMessage(song, channelID)
}

func NextSong() Song {
	fmt.Println("Next song...")
	if len(currentDJ.Queue) == 0 {
		return nil
	}

	if len(currentDJ.Queue) == 1 {
		currentDJ.CurrentSong = currentDJ.Queue[0]
		currentDJ.Queue = nil
		return currentDJ.CurrentSong
	}

	currentDJ.CurrentSong = currentDJ.Queue[0]
	currentDJ.Queue = currentDJ.Queue[1:]
	return currentDJ.CurrentSong
}
