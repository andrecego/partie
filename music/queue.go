package music

import "fmt"

func ShowQueue(channelID string) {
	queueMessage(channelID)
}

func Skip() {
	if currentDJ == nil || currentDJ.CurrentSong == nil {
		return
	}

	currentDJ.NeedsToSkip = true
}

func PlayPause() {
	if currentDJ == nil || currentDJ.CurrentSong == nil {
		return
	}

	if currentDJ.Paused {
		Resume()
	} else {
		Pause()
	}
}

func Pause() {
	currentDJ.Paused = true
}

func Resume() {
	currentDJ.Paused = false
}

func AddToQueue(song Song, channelID string) {
	currentDJ.Queue = append(currentDJ.Queue, song)
	updateQueueMessage()
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
