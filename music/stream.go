package music

import (
	"fmt"
	"io"
	"time"

	"github.com/andrecego/dca"
	"github.com/bwmarrin/discordgo"
)

func Stream(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if currentDJ.CurrentSong != nil {
		return nil
	}

	if currentDJ.Discord.VoiceConnection == nil {
		vs, err := FindVoiceChannel(s, m.GuildID, m.Author.ID)
		if err != nil {
			return fmt.Errorf("Error finding voice channel: %s", err)
		}

		err = Connect(vs.GuildID, vs.ChannelID)
		if err != nil {
			return fmt.Errorf("Error connecting to voice channel: %s", err)
		}
	}

	song := NextSong()
	updateQueueMessage()
	if song == nil {
		return fmt.Errorf("No song in queue")
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 180
	options.Application = "voip"
	options.Volume = 24
	// options.StartTime = 110

	// TODO: Add error check to see if url is still valid

	encodingSession, err := dca.EncodeFile(song.GetURL(), options)
	if err != nil {
		return fmt.Errorf("error encoding: %s", err)
	}
	defer encodingSession.Cleanup()

	done := make(chan error)
	streamSession := dca.NewStream(encodingSession, currentDJ.Discord.VoiceConnection, done)
	// tickerS := time.NewTicker(time.Second)
	tickerMs := time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case err = <-done:
			if err != nil && err != io.EOF {
				return fmt.Errorf("error streaming: %s", err)
			}

			if err == io.EOF {
				currentDJ.CurrentSong = nil
				Stream(s, m)
			}

		// case <-tickerS.C:
		// 	stats := encodingSession.Stats()
		// 	playbackPosition := streamSession.PlaybackPosition()

		// 	fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)

		case <-tickerMs.C:
			if currentDJ.Paused && !streamSession.Paused() {
				streamSession.SetPaused(true)
				continue
			}

			if streamSession.Paused() && !currentDJ.Paused {
				streamSession.SetPaused(false)
				continue
			}

			if currentDJ.NeedsToSkip {
				currentDJ.NeedsToSkip = false
				currentDJ.CurrentSong = nil
				encodingSession.Cleanup()
				Stream(s, m)
				continue
			}
		}
	}

}

func FindVoiceChannel(s *discordgo.Session, guildID, authorID string) (*discordgo.VoiceState, error) {
	g, err := s.State.Guild(guildID)
	if err != nil {
		fmt.Println("Error finding guild:", err)
		return nil, err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == authorID {
			return vs, nil
		}
	}

	return nil, fmt.Errorf("Could not find voice channel")
}
