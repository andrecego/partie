package voice_chat

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MoveUserBackAndForth(s *discordgo.Session, guildID, userID, otherChannelID string) error {
	// Check if the user is in a voice channel
	channelID := searchVoiceChannel(s, userID)
	if channelID == "" {
		return errors.New("User is not in a voice channel")
	}

	// Move to a different channel
	err := s.GuildMemberMove(guildID, userID, &otherChannelID)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Wait for the user to connect to the other channel and move back to the original channel
	for maxTimes := 1; maxTimes < 10; maxTimes++ {
		time.Sleep(10 * time.Millisecond)

		err = s.GuildMemberMove(guildID, userID, &channelID)
		if err == nil {
			break
		}

		fmt.Println(err.Error())
	}

	return nil
}

func searchVoiceChannel(session *discordgo.Session, user string) (voiceChannelID string) {
	for _, guild := range session.State.Guilds {
		for _, v := range guild.VoiceStates {
			if v.UserID == user {
				return v.ChannelID
			}
		}
	}

	return ""
}
