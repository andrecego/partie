package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

var ErrEmptyString = errors.New("Empty string")

var processedMessages = map[string]bool{}

var (
	messageCreateCommands = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, args []string){
		"music play": musicPlay,
	}
)

func musicPlay(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	err := addToQueue(s, strings.Join(args, " "), m.ChannelID, m.Author)
	if err != nil && !errors.Is(err, ErrEmptyString) {
		fmt.Println(err)
		return
	}

	err = music.Stream(s, m)
	if err != nil {
		fmt.Println(err)
		return
	}
}

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
	case "play":
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
	if query == "" {
		return fmt.Errorf("addToQueue: %w", ErrEmptyString)
	}

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

func PlaylistChannelHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != "955146633203560468" { // playlist channel
		return
	}

	if m.Author.ID == "943312702372192256" { // partie bot id
		return
	}

	go deleteMessage(s, m)

	if isCommand(m.Content) {
		return
	}

	music.New(s)

	musicPlay(s, m, []string{m.Content})
}

func deleteMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	time.Sleep(2500 * time.Millisecond)
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		fmt.Println("ERROR deleting message: ", err)
	}
}

func AddMusicReactionHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != "955146633203560468" { // playlist channel
		return
	}

	if m.Author.ID == "943312702372192256" { // partie bot id
		return
	}

	if !isCommand(m.Content) {
		return
	}

	command, args := parseCommands(m, "music")
	if command == "" || len(args) == 0 {
		return
	}

	if args[0] == "changeImage" {
		_, err := s.ChannelMessageEdit("955146633203560468", "955235579086389318", "https://preview.redd.it/xg4ke9fjvng41.jpg?width=1024&auto=webp&s=57153c7a8162153d2fd2a02b3d7bdc6085396c9f")
		if err != nil {
			fmt.Println(err)
		}

	}

	if args[0] != "addReact" {
		return
	}

	queueMessageID := "955235598292103199"
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "⏯️")
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "⏹️")
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "⏭️")
}

func ReactionControlHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.ChannelID != "955146633203560468" { // playlist channel
		return
	}

	if m.UserID == "943312702372192256" { // partie bot id
		return
	}

	switch m.Emoji.Name {
	case "⏯️":
		music.PlayPause()
	case "⏹️":
		// music.Stop()
	case "⏭️":
		music.Skip()
	}

	err := s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.Name, m.UserID)
	if err != nil {
		err = s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.ID, m.UserID)
		if err != nil {
			fmt.Println("ERROR removing reaction: ", err)
		}
	}
}
