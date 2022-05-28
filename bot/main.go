package bot

import (
	"context"
	"errors"
	"fmt" //to print errors
	"math/rand"
	"partie-bot/cache"
	"partie-bot/commands"
	"partie-bot/config" //importing our config package which we have created above
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo" //discordgo package from the repo of bwmarrin .
)

var BotId string
var goBot *discordgo.Session

const (
	afkChannel      = "950007673917698069"
	mutadinhoRoleID = "943703035115356191"
)

type VoiceStatusUpdate struct {
	UserID     string `json:"user_id"`
	ChannelID  string `json:"channel_id"`
	GuildID    string `json:"guild_id"`
	SessionID  string `json:"session_id"`
	Suppress   bool   `json:"suppress"`
	SelfVideo  bool   `json:"self_video"`
	SelfStream bool   `json:"self_stream"`
	SelfMute   bool   `json:"self_mute"`
	SelfDeaf   bool   `json:"self_deaf"`
	Mute       bool   `json:"mute"`
	Deaf       bool   `json:"deaf"`
}

func Start() {

	//creating new bot session
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	//Handling error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	//Handlinf error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	goBot.Identify.Intents = discordgo.IntentsAll

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(pingHandler)
	// goBot.AddHandler(notifyBadNameHandler) // roles ok, but not moving back to channel
	// goBot.AddHandler(allEventsHandler)
	goBot.AddHandler(subscribeToNameHandler)
	goBot.AddHandler(stopStreamHandler)
	goBot.AddHandler(addBlockedUserStreamHandler)
	goBot.AddHandler(removeBlockedUserStreamHandler)
	goBot.AddHandler(listBlockedUserStreamHandler)
	goBot.AddHandler(isStreamingHandler)
	goBot.AddHandler(voiceStateUpdateHandler)
	goBot.AddHandler(commands.RollD20Handler)
	goBot.AddHandler(commands.MusicHandler)
	goBot.AddHandler(commands.PlaylistChannelHandler)
	goBot.AddHandler(commands.PlaylistChannelStartHandler)
	goBot.AddHandler(commands.AddMusicReactionHandler)
	goBot.AddHandler(commands.ReactionControlHandler)
	// goBot.AddHandler(streamStartHandler) // infinite looping

	err = goBot.Open()
	//Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//If every thing works fine we will be printing this.
	fmt.Println("Bot is running !")
}

func msgWithPrefix(name string) string {
	return config.BotPrefix + name
}

func isStreamingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	contents := strings.Split(m.Content, " ")

	if contents[0] != msgWithPrefix("streaming?") {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!streaming? @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@!(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	userID := matches[1]

	vsu, err := s.State.VoiceState(m.GuildID, userID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting voice state: "+err.Error())
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The user is streaming? %v", vsu.SelfStream))
}

func addBlockedUserStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	contents := strings.Split(m.Content, " ")

	if contents[0] != msgWithPrefix("blockstream") {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!blockstream @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@!(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	userID := matches[1]

	err := cache.New().Client.SAdd(context.TODO(), stopStreamKey(m.GuildID), userID).Err()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error adding user to blocked list: "+err.Error())
		return
	}

	vsu, err := s.State.VoiceState(m.GuildID, userID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting voice state: "+err.Error())
		return
	}

	voiceStateUpdateHandler(s, &discordgo.VoiceStateUpdate{VoiceState: vsu})

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func removeBlockedUserStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	contents := strings.Split(m.Content, " ")

	if contents[0] != msgWithPrefix("rmblocked") {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!rmblocked @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@!(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	userID := matches[1]

	err := cache.New().Client.SRem(context.TODO(), stopStreamKey(m.GuildID), userID).Err()
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Error removing user of blocked list: "+err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func listBlockedUserStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content != msgWithPrefix("listblocked") {
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

func stopStreamKey(guildID string) string {
	return "guilds:" + guildID + ":blockStream"
}

func allEventsHandler(s *discordgo.Session, e *discordgo.Event) {
	fmt.Println()
	fmt.Println(e.Type, " - ", string(e.RawData))
}

func voiceStateUpdateHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.UserID == BotId || vsu.ChannelID == afkChannel {
		return
	}

	blockedIds, err := cache.New().Client.SMembers(context.TODO(), stopStreamKey(vsu.GuildID)).Result()
	if err != nil {
		fmt.Println("Error getting blocked users: " + err.Error())
		return
	}

	for _, id := range blockedIds {
		if id == vsu.UserID {
			if vsu.SelfStream == true {
				_, _ = s.ChannelMessageSend("943655307626823771", "Hey <@"+vsu.UserID+">, no stream for you.")
				moveUserBackAndForward(s, vsu.GuildID, vsu.UserID)
			}

			return
		}
	}
}

func handleVoiceStateUpdate(s *discordgo.Session, vsu VoiceStatusUpdate) {
	if vsu.UserID == BotId || vsu.ChannelID == afkChannel {
		return
	}

	blockedIds, err := cache.New().Client.SMembers(context.TODO(), stopStreamKey(vsu.GuildID)).Result()
	if err != nil {
		fmt.Println("Error getting blocked users: " + err.Error())
		return
	}

	for _, id := range blockedIds {
		if id == vsu.UserID {
			if vsu.SelfStream == true {
				_, _ = s.ChannelMessageSend("943655307626823771", "Hey <@"+vsu.UserID+">, no stream for you.")
				moveUserBackAndForward(s, vsu.GuildID, vsu.UserID)
			}

			return
		}
	}
}

func streamStartHandler(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	if vsu.ChannelID == "" {
		return
	}

	if vsu.UserID == BotId || vsu.UserID != "176049727945572352" {
		return
	}

	_, _ = s.ChannelMessageSend("943655307626823771", "Hey <@"+vsu.UserID+">, no stream for you.")
	moveUserBackAndForward(s, vsu.GuildID, vsu.UserID)
}

func subscribeToNameHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	if m.Content == msgWithPrefix("join") {
		guild, err := s.Guild(m.GuildID)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(guild.Name)
		fmt.Println(guild.ID)

		fmt.Println(searchVoiceChannel(s, m.Author.ID))
		_, _ = s.ChannelMessageSend(m.ChannelID, "Joined!")
	}
}

// SearchVoiceChannel search the voice channel id into from guild.
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

func stopStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "176049727945572352" {
		return
	}

	contents := strings.Split(m.Content, " ")

	if contents[0] != msgWithPrefix("bugoff") {
		return
	}

	if len(contents) != 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid option: should be `!bugoff @user`")
		return
	}

	userIdRegex := regexp.MustCompile(`<@!(\d*)>`)
	matches := userIdRegex.FindStringSubmatch(contents[1])
	userID := matches[1]

	err := moveUserBackAndForward(s, m.GuildID, userID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func moveUserBackAndForward(s *discordgo.Session, guildID, userID string) error {
	// Check if the user is in a voice channel
	channelID := searchVoiceChannel(s, userID)
	if channelID == "" {
		return errors.New("User is not in a voice channel")
	}

	// Move to a different channel
	afkChanID := afkChannel
	err := s.GuildMemberMove(guildID, userID, &afkChanID)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Wait for the user to connect to the AFK channel and move back to the original channel
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

//Definition of pingHandler function it takes two arguments first one is discordgo.Session which is s , second one is discordgo.MessageCreate which is m.
func pingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotId {
		return
	}

	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if m.Content == "whoami" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">")
	}
}

func notifyBadNameHandler(s *discordgo.Session, gm *discordgo.GuildMemberUpdate) {
	if gm.User.ID == BotId {
		return
	}

	var name string
	if len(gm.Nick) == 0 {
		name = gm.User.Username
	} else {
		name = gm.Nick
	}

	if strings.ToLower(name[0:1]) == "e" {
		err := s.GuildMemberRoleRemove(gm.GuildID, gm.User.ID, mutadinhoRoleID)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	_, err := s.ChannelMessageSend("943655307626823771", "<@"+gm.User.ID+"> esse nick `"+name+"` ta errado ein")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = s.GuildMemberRoleAdd(gm.GuildID, gm.User.ID, mutadinhoRoleID)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = moveUserBackAndForward(s, gm.GuildID, gm.User.ID)
	if err != nil {
		fmt.Println(err.Error())
	}

}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func changeVoiceName(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotId {
		return
	}

	if m.Content == "newDay" {
		rand.Seed(time.Now().UnixNano())

		fmt.Println("Changing voice name")

		letterOfTheDay := fmt.Sprintf("[%s] - Letra do dia ", randSeq(1))

		_, err := s.ChannelEdit("833404286109089822", letterOfTheDay)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
