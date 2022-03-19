package music

import (
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func Stream(s *discordgo.Session, m *discordgo.MessageCreate) error {
	vs, err := FindVoiceChannel(s, m.GuildID, m.Author.ID)
	if err != nil {
		return fmt.Errorf("Error finding voice channel: %s", err)
	}

	err = Connect(vs.GuildID, vs.ChannelID)
	if err != nil {
		return fmt.Errorf("Error connecting to voice channel: %s", err)
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 256
	options.Application = "voip"

	downloadURL := "https://rr6---sn-8p8v-bg0ls.googlevideo.com/videoplayback?expire=1647726526&ei=Xvs1YpTZJ66y1sQPq82x0A8&ip=187.35.3.130&id=o-AEvD0vgkKQCbEZ0tby49fq-twPcZCf_J4YpTd9AbDJyV&itag=251&source=youtube&requiressl=yes&mh=s3&mm=31%2C29&mn=sn-8p8v-bg0ls%2Csn-b8u-bg0e&ms=au%2Crdu&mv=m&mvi=6&pl=24&gcr=br&initcwndbps=1037500&vprv=1&mime=audio%2Fwebm&gir=yes&clen=3104544&dur=189.841&lmt=1632589970085732&mt=1647704447&fvip=6&keepalive=yes&fexp=24001373%2C24007246&c=ANDROID&txp=5511222&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cgcr%2Cvprv%2Cmime%2Cgir%2Cclen%2Cdur%2Clmt&sig=AOq0QJ8wRAIgIALJQ_9SHKVxKauNS9biC48ycwhKiPCLC5Zgcds2TtUCIH9-Us5-dDY4CFIwMztTDusztn58PnB-smdWJn01YAMw&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRQIgNp7fX3w_6IZ2yEAruob-wV0-VwCnu18ndLPmVEVIDCsCIQD4boEpJcEC20NsK1BOPWoL9ieg0nm-MZqOqERbTRgp6A%3D%3D"

	encodingSession, err := dca.EncodeFile(downloadURL, options)
	if err != nil {
		return fmt.Errorf("error encoding: %s", err)
	}
	defer encodingSession.Cleanup()

	done := make(chan error)
	streamSession := dca.NewStream(encodingSession, currentDJ.Discord.VoiceConnection, done)

	for {
		select {
		case err = <-done:
			fmt.Print("Done")
			fmt.Println(err)
			if err != nil && err != io.EOF {
				return fmt.Errorf("error streaming: %s", err)
			}

			finished, err := streamSession.Finished()
			if err != nil {
				return fmt.Errorf("error checking if stream finished: %s", err)
			}

			if finished {
				return nil
			}

		default:
			time.Sleep(50 * time.Millisecond)

			if currentDJ.Paused {
				streamSession.SetPaused(true)
				continue
			}

			if streamSession.Paused() && !currentDJ.Paused {
				streamSession.SetPaused(false)
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
