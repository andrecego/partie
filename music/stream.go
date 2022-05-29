package music

import (
	"fmt"
	"io"
	"time"

	"github.com/andrecego/dca"
	"github.com/bwmarrin/discordgo"
)

func Stream(s *discordgo.Session, message *discordgo.MessageCreate) error {
	if currentDJ.CurrentSong != nil {
		return nil
	}

	if currentDJ.Discord.VoiceConnection == nil {
		vs, err := FindVoiceChannel(s, message.GuildID, message.Author.ID)
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
	volume := -40 // -70..-5 where -5 is max volume
	// loudnorm -> Indicates the name of the normalization filter
	// I, i     -> Indicates the integrated loudness (-70 to -5.0 with default -24.0)
	// LRA, lra -> Indicates the loudness range (1.0 to 20.0 with default 7.0)
	// TP, tp   -> Indicates the max true peak (-9.0 to 0.0 with default -2.0)
	options.AudioFilter = fmt.Sprintf("loudnorm=I=%v:LRA=11:TP=-5", volume)
	// options.StartTime = 110
	// options.Volume = 8 // cant be used with audio filter

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
			if err == nil {
				continue
			}

			if err == io.EOF || err == io.ErrUnexpectedEOF {
				if err == io.ErrUnexpectedEOF {
					fmt.Printf("Skipping song `%v` due to unexpected EOF", song.GetTitle())
				}

				currentDJ.CurrentSong = nil
				Stream(s, message)
				return nil
			}

			return fmt.Errorf("error streaming: %s", err)

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
				Stream(s, message)
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
