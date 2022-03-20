package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"partie-bot/cache"
	"partie-bot/music"
	"partie-bot/music/youtube"
	"partie-bot/repositories"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var buffer = make([][]byte, 0)

func getSearchResult(guildID, userID string) ([]repositories.YoutubeSearchResult, error) {
	key := fmt.Sprintf("guilds:%s:users:%s:search_results", guildID, userID)

	bs, err := cache.New().Client.Get(context.TODO(), key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("Error getting search result: %s", err)
	}

	var results []repositories.YoutubeSearchResult
	err = json.Unmarshal(bs, &results)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling search result: %s", err)
	}

	return results, nil
}

func MusicHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// searchResults, err := getSearchResult(m.GuildID, m.Author.ID)
	// if err != nil {
	// 	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error getting search results: %s", err))
	// 	return
	// }

	if !isCommand(m.Content) {
		return
	}

	command, args := parseCommands(m, "music")
	if command != "music" {
		return
	}

	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !music <play/pause/stop>")
		return
	}

	music.New(s)

	switch args[0] {
	// case "play":
	// 	s.ChannelMessageSend(m.ChannelID, "Playing music")
	// 	err := play(s, strings.Join(args[1:], " "), m.GuildID, m.Author.ID)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	case "play":
		// s.ChannelMessageSend(m.ChannelID, "Playing music")
		s.MessageReactionAdd(m.ChannelID, m.Message.ID, "✅")

		query := strings.Join(args[1:], " ")
		if query != "" {
			err := addToQueue(s, strings.Join(args[1:], " "), m.ChannelID, m.Author)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err := music.Stream(s, m)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "add":
		s.ChannelMessageSend(m.ChannelID, "Adding music")
		err := addToQueue(s, strings.Join(args[1:], " "), m.ChannelID, m.Author)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "skip":
		music.Skip()
	case "search":
		music.Search(strings.Join(args[1:], " "), m)
	case "queue":
		music.ShowQueue(m.ChannelID)
	case "stream":
		err := music.Stream(s, m)
		if err != nil {
			fmt.Println(err)
			return
		}

	case "pause":
		music.Pause()
		s.ChannelMessageSend(m.ChannelID, "Pausing music")
	case "resume":
		music.Resume()
		s.ChannelMessageSend(m.ChannelID, "Resuming music")
	case "stop":
		s.ChannelMessageSend(m.ChannelID, "Stopping music")
	default:
		s.ChannelMessageSend(m.ChannelID, "Usage: !music <play/pause/stop>")
	}
}

func addToQueue(s *discordgo.Session, query, channelID string, author *discordgo.User) error {
	finder := music.ParseQuery(query)

	jsonInfo, err := finder.Download()
	if err != nil {
		return fmt.Errorf("Error downloading video: %s", err)
	}

	var song youtube.Youtube
	err = json.Unmarshal([]byte(jsonInfo), &song)
	if err != nil {
		return fmt.Errorf("Error unmarshalling info file: %s", err)
	}
	song.AddedBy = author

	music.AddToQueue(&song, channelID)

	return nil
}

func shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	fmt.Println("Running command:", command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
