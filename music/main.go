package music

import (
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
	delayToSpeak      = 250 * time.Millisecond
	timesToDisconnect = 60 * 5
)

func New(s *discordgo.Session) *DJ {
	if currentDJ != nil {
		return currentDJ
	}

	currentDJ = &DJ{
		NeedsToSkip: false,
		Discord: &Discord{
			Session: s,
		},
	}

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
