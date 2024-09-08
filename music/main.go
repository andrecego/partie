package music

import (
	"context"
	"encoding/json"
	"fmt"
	"partie-bot/cache"
	"partie-bot/music/youtube"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DJ struct {
	CurrentSong Song
	Queue       []Song
	Discord     *Discord
	Volume      float64
	Paused      bool
	NeedsToSkip bool
}

type Discord struct {
	Session         *discordgo.Session
	VoiceConnection *discordgo.VoiceConnection
	VoiceState      *discordgo.VoiceState
}

type Song interface {
	GetID() string
	GetTitle() string
	GetDuration() time.Duration
	GetAddedBy() string
	GetVideoURL() string
	GetURL() string
	GetThumbnail() string
	GetStartTime() int
	GetGuildID() string
	GetChannelID() string
	GetAuthorID() string
}

var (
	currentDJ         *DJ
	delayToSpeak      = 50 * time.Millisecond
	timesToDisconnect = 60 * 5
)

func New(s *discordgo.Session, guildID string) *DJ {
	if currentDJ != nil {
		return currentDJ
	}

	redisClient := cache.New().Client
	key := "guilds:" + guildID + ":queue"
	allSongsBytes, err := redisClient.Get(context.TODO(), key).Bytes()
	if err != nil {
		fmt.Println("Error getting queue from cache:", err)
	}

	var allYoutubeSongs []youtube.Youtube
	if allSongsBytes != nil {
		err = json.Unmarshal(allSongsBytes, &allYoutubeSongs)
		if err != nil {
			fmt.Println("allSongstext:", string(allSongsBytes))
			fmt.Println("Error unmarshalling queue:", err)
		}
	}

	var allSongs []Song
	for _, song := range allYoutubeSongs {
		allSongs = append(allSongs, &song)
	}

	currentDJ = &DJ{
		NeedsToSkip: false,
		Queue:       allSongs,
		CurrentSong: nil,
		Discord: &Discord{
			Session: s,
		},
	}

	updateQueueMessage()

	return currentDJ
}

func Connect(guildID, channelID string) error {
	if currentDJ.Discord.VoiceConnection != nil {
		return nil
	}

	// Join the provided voice channel.
	// guildID := currentDJ.Discord.VoiceState.GuildID
	// channelID := currentDJ.Discord.VoiceState.ChannelID
	vc, err := currentDJ.Discord.Session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	currentDJ.Discord.VoiceConnection = vc

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(delayToSpeak)

	// Start speaking.
	currentDJ.Discord.VoiceConnection.Speaking(true)

	return nil
}

func SetVoiceState(vs *discordgo.VoiceState) {
	if currentDJ.Discord.VoiceState != nil {
		return
	}

	currentDJ.Discord.VoiceState = vs
}
