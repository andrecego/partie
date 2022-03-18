package music

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DJ struct {
	CurrentSong Song
	Queue       []Song
	Buffer      [][]byte
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
	GetAddedBy() string
	GetPath() string
}

var (
	currentDJ         *DJ
	delayToSpeak      = 250 * time.Millisecond
	timesToDisconnect = 60 * 5
)

func New(s *discordgo.Session) *DJ {
	if currentDJ != nil {
		return currentDJ
	}

	currentDJ = &DJ{
		Buffer:      make([][]byte, 0),
		NeedsToSkip: false,
		Discord: &Discord{
			Session: s,
		},
	}

	return currentDJ
}

func Connect() error {
	if currentDJ.Discord.VoiceConnection != nil {
		return nil
	}

	// Join the provided voice channel.
	guildID := currentDJ.Discord.VoiceState.GuildID
	channelID := currentDJ.Discord.VoiceState.ChannelID
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

func Disconnect() {
	// Sleep for a specificed amount of time before ending.
	for i := 0; i < timesToDisconnect; i++ {
		if len(currentDJ.Queue) > 0 {
			Play()
			return
		}

		time.Sleep(time.Second)
	}

	// Stop speaking
	currentDJ.Discord.VoiceConnection.Speaking(false)

	// Disconnect from the provided voice channel.
	currentDJ.Discord.VoiceConnection.Disconnect()
}

func Play() error {
	if currentDJ.CurrentSong != nil {
		return nil
	}

	err := Connect()
	if err != nil {
		return fmt.Errorf("Error connecting to voice channel: %s", err)
	}

	NextSong()
	for currentDJ.CurrentSong != nil {
		fmt.Println("Flushing buffer...")
		currentDJ.Buffer = make([][]byte, 0)

		err := addToBuffer(currentDJ.CurrentSong)
		if err != nil {
			return fmt.Errorf("Error adding to buffer: %s", err)
		}

		nowPlayingMessage()
		PlayLoop()
	}

	// Disconnect()

	return nil
}

func PlayLoop() {
	for _, buff := range currentDJ.Buffer {
		if currentDJ.NeedsToSkip {
			skipSong()
			return
		}

		currentDJ.Discord.VoiceConnection.OpusSend <- buff
	}

	skipSong()
}
