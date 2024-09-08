package commands

import (
	"context"
	"fmt"
	"partie-bot/cache"
	"partie-bot/config"
	"partie-bot/helpers"
	"partie-bot/voice_chat"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func AddStreamBlockCommands(goBot *discordgo.Session) {
	goBot.AddHandler(addBlockedUserStreamHandler)
	goBot.AddHandler(removeBlockedUserStreamHandler)
	goBot.AddHandler(listBlockedUserStreamHandler)
	goBot.AddHandler(blockStreamHandler)
	goBot.AddHandler(stopStreamHandler)

	goBot.AddHandler(addBlockedUserVideoHandler)
	goBot.AddHandler(removeBlockedUserVideoHandler)
	goBot.AddHandler(listBlockedUserVideoHandler)
	goBot.AddHandler(blockVideoHandler)
	goBot.AddHandler(stopVideoHandler)
}

func addBlockedUserStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	command, args := helpers.ParseCommand(m.Content, config.DogeGuildConfig)
	if command != "blockstream" {
		return
	}

	if len(args) != 1 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!blockstream @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(args[0])
	if len(matches) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid user format")
		return
	}
	userID := matches[1]

	err := cache.New().Client.SAdd(context.TODO(), stopStreamKey(m.GuildID), userID).Err()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error adding user to blocked list: "+err.Error())
		return
	}

	vsu, err := s.State.VoiceState(m.GuildID, userID)
	if err != nil && err != discordgo.ErrStateNotFound {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting voice state: "+err.Error())
		return
	}
	if err != discordgo.ErrStateNotFound {
		blockStreamHandler(s, &discordgo.VoiceStateUpdate{VoiceState: vsu})
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func addBlockedUserVideoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	contents := strings.Split(m.Content, " ")

	if contents[0] != helpers.MsgWithPrefix("blockvideo", config.DogeGuildConfig) {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!blockvideo @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	if len(matches) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid user format")
		return
	}
	userID := matches[1]

	err := cache.New().Client.SAdd(context.TODO(), stopVideoKey(m.GuildID), userID).Err()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error adding user to blocked list: "+err.Error())
		return
	}

	vsu, err := s.State.VoiceState(m.GuildID, userID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting voice state: "+err.Error())
		return
	}

	blockStreamHandler(s, &discordgo.VoiceStateUpdate{VoiceState: vsu})

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}
func removeBlockedUserStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	contents := strings.Split(m.Content, " ")

	if contents[0] != helpers.MsgWithPrefix("rmBlockedStream", config.DogeGuildConfig) {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!rmblocked @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	if len(matches) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid user format")
		return
	}
	userID := matches[1]

	err := cache.New().Client.SRem(context.TODO(), stopStreamKey(m.GuildID), userID).Err()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error removing user of blocked list: "+err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func removeBlockedUserVideoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	contents := strings.Split(m.Content, " ")

	if contents[0] != helpers.MsgWithPrefix("rmBlockedVideo", config.DogeGuildConfig) {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!rmblocked @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	if len(matches) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid user format")
		return
	}
	userID := matches[1]

	err := cache.New().Client.SRem(context.TODO(), stopVideoKey(m.GuildID), userID).Err()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error removing user of blocked list: "+err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func listBlockedUserStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content != helpers.MsgWithPrefix("listBlockedStream", config.DogeGuildConfig) {
		return
	}

	blockedIds, err := cache.New().Client.SMembers(context.TODO(), stopStreamKey(m.GuildID)).Result()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error listing blocked users: "+err.Error())
		return
	}

	if len(blockedIds) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No blocked users")
		return
	}

	message := "Users blocked:\n"
	for _, id := range blockedIds {
		message = message + "<@!" + id + ">" + "\n"
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}

func listBlockedUserVideoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content != helpers.MsgWithPrefix("listBlockedVideo", config.DogeGuildConfig) {
		return
	}

	blockedIds, err := cache.New().Client.SMembers(context.TODO(), stopVideoKey(m.GuildID)).Result()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error listing blocked users: "+err.Error())
		return
	}

	if len(blockedIds) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No blocked users")
		return
	}

	message := "Users blocked:\n"
	for _, id := range blockedIds {
		message = message + "<@!" + id + ">" + "\n"
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, message)
}

func blockStreamHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.UserID == config.BotId {
		return
	}

	if vsu.SelfStream == false {
		return
	}

	blockedIds, err := cache.New().Client.SMembers(context.TODO(), stopStreamKey(vsu.GuildID)).Result()
	if err != nil {
		fmt.Println("Error getting blocked users: " + err.Error())
		return
	}

	blockedUser := slices.Contains(blockedIds, vsu.UserID)
	if blockedUser {
		_, _ = s.ChannelMessageSend("943655307626823771", "Hey <@"+vsu.UserID+">, no stream for you.")
		voice_chat.MoveUserBackAndForth(s, vsu.GuildID, vsu.UserID, config.DogeGuildConfig.AfkChannelId)
	}
}

func blockVideoHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.UserID == config.BotId || vsu.ChannelID == config.DogeGuildConfig.AfkChannelId {
		return
	}

	blockedIds, err := cache.New().Client.SMembers(context.TODO(), stopVideoKey(vsu.GuildID)).Result()
	if err != nil {
		fmt.Println("Error getting blocked users: " + err.Error())
		return
	}

	for _, id := range blockedIds {
		if id == vsu.UserID {
			if vsu.SelfVideo == true {
				_, _ = s.ChannelMessageSend("943655307626823771", "Hey <@"+vsu.UserID+">, no video for you.")
				voice_chat.MoveUserBackAndForth(s, vsu.GuildID, vsu.UserID, config.DogeGuildConfig.AfkChannelId)
			}

			return
		}
	}
}

func stopStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "176049727945572352" {
		return
	}

	contents := strings.Split(m.Content, " ")

	if contents[0] != helpers.MsgWithPrefix("stopstream", config.DogeGuildConfig) {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!stopstream @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	if len(matches) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid user format")
		return
	}
	userID := matches[1]

	err := voice_chat.MoveUserBackAndForth(s, m.GuildID, userID, config.DogeGuildConfig.AfkChannelId)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func stopVideoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "176049727945572352" {
		return
	}

	contents := strings.Split(m.Content, " ")

	if contents[0] != helpers.MsgWithPrefix("stopvideo", config.DogeGuildConfig) {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!stopvideo @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	if len(matches) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid user format")
		return
	}
	userID := matches[1]

	err := voice_chat.MoveUserBackAndForth(s, m.GuildID, userID, config.DogeGuildConfig.AfkChannelId)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func stopStreamKey(guildID string) string {
	return "guilds:" + guildID + ":blockStream"
}

func stopVideoKey(guildID string) string {
	return "guilds:" + guildID + ":blockVideo"
}
