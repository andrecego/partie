package music

import (
	"context"
	"encoding/json"
	"fmt"
	"partie-bot/cache"

	"github.com/bwmarrin/discordgo"
)

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

func Cleanup(session *discordgo.Session, guildID string) {
	currentDJ.CurrentSong = nil
	currentDJ.Queue = nil
	updateQueueCache()
	currentDJ = nil

	New(session, guildID)
}

func AddAsyncToQueue(song Song) {
	addToQueue(song)
}

func AddToQueue(song Song) {
	addToQueue(song)
	updateQueueMessage()
	addedToQueueMessage(song)
}

func addToQueue(song Song) {
	currentDJ.Queue = append(currentDJ.Queue, song)
	updateQueueCache()
}

func updateQueueCache() {
	redisClient := cache.New().Client
	key := fmt.Sprintf("guilds:%s:queue", currentDJ.GuildID)
	allSongs := currentDJ.Queue
	if currentDJ.CurrentSong != nil {
		allSongs = append([]Song{currentDJ.CurrentSong}, allSongs...)
	}

	if len(allSongs) == 0 {
		err := redisClient.Del(context.TODO(), key).Err()
		if err != nil {
			fmt.Println("Error deleting queue: ", err)
		}
		return
	}

	allSongsBytes, err := json.Marshal(allSongs)
	if err != nil {
		fmt.Println("Error marshalling queue: ", err)
		return
	}

	err = redisClient.Set(context.TODO(), key, allSongsBytes, 0).Err()
	if err != nil {
		fmt.Println("Error saving queue: ", err)
	}
}

func UpdateQueueMessage() {
	updateQueueMessage()
}

func Remove(queueNumber int) {
	fmt.Println("Removing song from queue: ", queueNumber)
	index := queueNumber - 1
	if index < 0 || index >= len(currentDJ.Queue) {
		return
	}

	currentDJ.Queue = append(currentDJ.Queue[:index], currentDJ.Queue[index+1:]...)
	updateQueueCache()
	updateQueueMessage()
}

func NextSong() Song {
	nextSong := fetchNextSong()
	if nextSong != nil {
		fmt.Println("Next song: ", nextSong.GetTitle())
	}

	go updateQueueCache()
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
