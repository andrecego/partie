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

func AddToQueue(song Song) {
	currentDJ.Queue = append(currentDJ.Queue, song)
	updateQueueMessage()
	addedToQueueMessage(song)
}

func Remove(queueNumber int) {
	fmt.Println("Removing song from queue: ", queueNumber)
	index := queueNumber - 1
	if index < 0 || index >= len(currentDJ.Queue) {
		return
	}

	currentDJ.Queue = append(currentDJ.Queue[:index], currentDJ.Queue[index+1:]...)
	updateQueueMessage()
}

func NextSong() Song {
	nextSong := fetchNextSong()
	fmt.Println("Next song...")

	go updateQueueMessage()
	return nextSong
}

func fetchNextSong() Song {
	if len(currentDJ.Queue) == 0 {
		currentDJ.Queue = nil
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
