package music

import (
	"context"
	"fmt"
	"io"
	"partie-bot/cache"
	"time"

	"github.com/andrecego/dca"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
)

const (
	disconnectAfterDuration = 3 * time.Minute
)

func keepAliveCacheKey(guildID string) string {
	return "guilds:" + guildID + ":keep_alive"
}

func setKeepAliveCache(guildID string, expires time.Duration) error {
	key := keepAliveCacheKey(guildID)
	err := cache.New().Client.Set(context.TODO(), key, true, expires).Err()
	if err != nil {
		fmt.Println("Error setting disconnect cache:", err)
	}

	return err
}

func Disconnect(session *discordgo.Session, guildID string) {
	if currentDJ == nil || currentDJ.Discord.VoiceConnection == nil {
		return
	}

	err := currentDJ.Discord.VoiceConnection.Disconnect()
	if err != nil {
		fmt.Println("Error disconnecting:", err)
	}

	Cleanup(session, guildID)
}

func disconnectAfter(session *discordgo.Session, guildID string, duration time.Duration) {
	err := setKeepAliveCache(guildID, duration)
	if err != nil {
		fmt.Println("Error setting keep alive cache:", err)
		return
	}

	time.AfterFunc(duration+2*time.Second, func() {
		key := keepAliveCacheKey(guildID)
		_, err := cache.New().Client.Get(context.TODO(), key).Result()
		if err == redis.Nil {
			Disconnect(session, guildID)
		} else if err != nil {
			fmt.Println("Error getting keep alive cache:", err)
		}
	})
}

func Stream(session *discordgo.Session) error {
	if currentDJ.CurrentSong != nil {
		return nil
	}

	currentSong := NextSong()
	if currentSong == nil {
		fmt.Println("No song in queue")
		return nil
	}

	if currentDJ.Discord.VoiceConnection == nil {
		vs, err := FindVoiceChannel(session, currentSong)
		if err != nil {
			return fmt.Errorf("Error finding voice channel: %s", err)
		}

		err = Connect(vs.GuildID, vs.ChannelID)
		if err != nil {
			return fmt.Errorf("Error connecting to voice channel: %s", err)
		}
	}

	disconnectAfterSongDuration := currentSong.GetDuration() + disconnectAfterDuration
	fmt.Printf("Setting keep alive cache for %.0f seconds\n", disconnectAfterSongDuration.Seconds())
	go disconnectAfter(session, currentDJ.Discord.VoiceConnection.GuildID, disconnectAfterSongDuration)

	// If there are more songs in the queue, start downloading the next song
	if len(currentDJ.Queue) > 0 {
		go currentDJ.Queue[0].GetURL()
	}

	encodingSession, err := dca.EncodeFile(currentSong.GetURL(), options(0))
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

			if err == io.EOF {
				currentDJ.CurrentSong = nil
				Stream(session)
				return nil
			}

			return fmt.Errorf("error streaming: %s", err)

		// case <-tickerS.C:
		// 	stats := encodingSession.Stats()
		// 	playbackPosition := streamSession.PlaybackPosition()

		// 	fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)

		case <-tickerMs.C:
			if currentDJ == nil || currentDJ.CurrentSong == nil {
				return nil
			}

			if currentDJ.Paused && !streamSession.Paused() {
				streamSession.SetPaused(true)
				continue
			}

			if streamSession.Paused() && !currentDJ.Paused {
				streamSession.SetPaused(false)
				continue
			}

			if currentDJ.NeedsToSkip {
				streamSession.Finished()
				currentDJ.NeedsToSkip = false
				currentDJ.CurrentSong = nil
				encodingSession.Cleanup()
				Stream(session)
				continue
			}
		}
	}
}

func options(startTime int) *dca.EncodeOptions {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 180
	options.Application = "voip"
	volume := -30 // -70..-5 where -5 is max volume
	// loudnorm -> Indicates the name of the normalization filter
	// I, i     -> Indicates the integrated loudness (-70 to -5.0 with default -24.0)
	// LRA, lra -> Indicates the loudness range (1.0 to 20.0 with default 7.0)
	// TP, tp   -> Indicates the max true peak (-9.0 to 0.0 with default -2.0)
	options.AudioFilter = fmt.Sprintf("loudnorm=I=%v:LRA=11:TP=-5", volume)
	options.StartTime = startTime // in seconds
	// options.Volume = 8 // cant be used with audio filter

	// TODO: Add error check to see if url is still valid

	return options
}

func FindVoiceChannel(s *discordgo.Session, song Song) (*discordgo.VoiceState, error) {
	g, err := s.State.Guild(song.GetGuildID())
	if err != nil {
		fmt.Println("Error finding guild:", err)
		return nil, err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == song.GetAuthorID() {
			return vs, nil
		}
	}

	return nil, fmt.Errorf("Could not find voice channel")
}
