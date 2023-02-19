package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"partie-bot/cache"
	"partie-bot/music"
	"partie-bot/music/youtube"
	"partie-bot/repositories"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rollbar/rollbar-go"
)

var buffer = make([][]byte, 0)

var ErrEmptyString = errors.New("Empty string")

var processedMessages = map[string]bool{}

// var (
// 	messageCreateCommands = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate, args []string){
// 		"music play": musicPlay,
// 	}
// )

func addedBy(message *discordgo.MessageCreate) youtube.AddedBy {
	return youtube.AddedBy{
		User:    message.Author,
		Guild:   &discordgo.Guild{ID: message.GuildID},
		Channel: &discordgo.Channel{ID: message.ChannelID},
	}
}

func musicPlay(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {

	err := addToQueue(session, strings.Join(args, " "), addedBy(message))
	if err != nil && !errors.Is(err, ErrEmptyString) {
		fmt.Println(err)
		return
	}

	err = music.Stream(session)
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
			err := addToQueue(s, strings.Join(args[1:], " "), addedBy(m))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		err := music.Stream(s)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "add":
		s.ChannelMessageSend(m.ChannelID, "Adding music")
		err := addToQueue(s, strings.Join(args[1:], " "), addedBy(m))
		if err != nil {
			fmt.Println(err)
			return
		}
	case "skip":
		music.Skip()
	case "search":
		music.Search(strings.Join(args[1:], " "), m)
	case "restart":
		music.Restart(s, m.GuildID, m.Author.ID)
	case "queue":
		music.ShowQueue(m.ChannelID)
	case "stream":
		err := music.Stream(s)
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
	case "createPlaylistChannel":
		s.ChannelMessageSend("955146633203560468", "Creating playlist channel")
		s.ChannelMessageSend("955146633203560468", "Creating playlist channel")
	default:
		s.ChannelMessageSend(m.ChannelID, "Usage: !music <play/pause/stop>")
	}
}

func addToQueue(s *discordgo.Session, query string, addedBy youtube.AddedBy) error {
	if query == "" {
		return fmt.Errorf("addToQueue: %w", ErrEmptyString)
	}

	queries := strings.Split(query, "\n")

	for i := range queries {
		finder := music.ParseQuery(queries[i])

		jsonInfo, err := finder.Download()
		if err != nil {
			return fmt.Errorf("Error downloading video: %s", err)
		}

		var youtubeResult youtube.YoutubeResult
		err = json.Unmarshal([]byte(jsonInfo), &youtubeResult)
		if err != nil {
			return fmt.Errorf("Error unmarshalling info file: %s", err)
		}

		if len(youtubeResult.Entries) == 0 {
			var youtubeEntry youtube.Youtube
			err = json.Unmarshal([]byte(jsonInfo), &youtubeEntry)
			if err != nil {
				return fmt.Errorf("Error unmarshalling info file: %s", err)
			}

			youtubeResult.Entries = append(youtubeResult.Entries, youtubeEntry)
		}

		for i := range youtubeResult.Entries {
			fmt.Println("Item: ", i)
			fmt.Println("Adding song: ", youtubeResult.Entries[i].Title)

			rollbar.SetPerson(addedBy.User.ID, addedBy.User.Username, "")
			rollbar.Info("Song added", map[string]interface{}{
				"title":    youtubeResult.Entries[i].Title,
				"videoURL": youtubeResult.Entries[i].VideoURL,
			})
			rollbar.ClearPerson()

			youtubeResult.Entries[i].AddedBy = addedBy
			music.AddAsyncToQueue(&youtubeResult.Entries[i])
		}
		music.UpdateQueueMessage()
	}

	return nil
}

func PlaylistChannelHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.ChannelID != "955146633203560468" { // playlist channel
		return
	}

	if message.Author.ID == "985002886087999608" { // partie bot id
		return
	}

	go deleteMessage(session, message)

	if isCommand(message.Content) {
		return
	}

	music.New(session)
	if isPrefixlessCommands(message.Content) {
		command, args := commandParse(message.Content)
		switch command {
		case "pause", "play":
			music.PlayPause()
		case "skip":
			music.Skip()
		case "restart":
			music.Restart(session, message.GuildID, message.Author.ID)
		case "remove", "delete":
			if len(args) == 0 {
				return
			}

			queueNumber, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Error converting queue number to int: ", err)
				return
			}

			music.Remove(queueNumber)
		}
		return
	}

	musicPlay(session, message, []string{message.Content})
}

func PlaylistChannelStartHandler(s *discordgo.Session, guild *discordgo.GuildCreate) {

	dogeID := "176049854001315850"
	if guild.Guild.ID != dogeID {
		return
	}

	playlistChannelID := "955146633203560468"
	messages, err := s.ChannelMessages(playlistChannelID, 100, "", "1052568412003520522", "")
	if err != nil {
		fmt.Println("Error getting messages:", err)
		return
	}

	for _, message := range messages {
		fmt.Println(message.Content)
		if message.ID == "1052568410799755264" || message.ID == "1052568412003520522" {
			continue
		}

		err := s.ChannelMessageDelete(playlistChannelID, message.ID)
		if err != nil {
			fmt.Println("Error deleting message:", err)
		}
	}

}

func deleteMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	time.Sleep(300 * time.Millisecond)
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		fmt.Println("ERROR deleting message: ", err)
	}
}

func AddMusicReactionHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != "955146633203560468" { // playlist channel
		return
	}

	if m.Author.ID == "985002886087999608" { // partie bot id
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
		_, err := s.ChannelMessageEdit("955146633203560468", "1052568410799755264", "https://preview.redd.it/xg4ke9fjvng41.jpg?width=1024&auto=webp&s=57153c7a8162153d2fd2a02b3d7bdc6085396c9f")
		if err != nil {
			fmt.Println(err)
		}

	}

	if args[0] != "addReact" {
		return
	}

	queueMessageID := "1052568412003520522"
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "⏯️")
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "⏹️")
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "⏭️")
	s.MessageReactionAdd(m.ChannelID, queueMessageID, "🔁")
}

func ReactionControlHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.ChannelID != "955146633203560468" { // playlist channel
		return
	}

	if m.UserID == "985002886087999608" { // partie bot id
		return
	}

	switch m.Emoji.Name {
	case "⏯️":
		music.PlayPause()
	case "⏹️":
		// music.Stop()
	case "⏭️":
		music.Skip()
	case "🔁":
		music.Restart(s, m.GuildID, m.UserID)
	}

	err := s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.Name, m.UserID)
	if err != nil {
		err = s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.ID, m.UserID)
		if err != nil {
			fmt.Println("ERROR removing reaction: ", err)
		}
	}
}

// DisconnectedHandler is called when the bot is disconnected from the server
func DisconnectedHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	// check if the user updated is the bot
	if vsu.UserID != s.State.User.ID {
		return
	}

	// check if the bot is being disconnected from the server
	if vsu.ChannelID != "" {
		return
	}

	// call the music cleanup function
	music.Cleanup()
	fmt.Println("Disconnected from server, music cleanup called")
}
