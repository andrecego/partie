package config

import "os"

var (
	Token           string // To store value of Token from config.json .
	RollbarToken    string
	BotId           string
	DogeGuildConfig = GuildConfig{
		GuildId:                "176049854001315850",
		Prefix:                 "!",
		PlaylistChannelId:      "1086983380828176505",
		PlaylistMessageImageId: "1222897849285742652",
		PlaylistMessageQueueId: "1222897838493532200",
		AfkChannelId:           "950007673917698069",
	}
)

func ReadConfig() error {
	Token = os.Getenv("TOKEN")
	RollbarToken = os.Getenv("ROLLBAR_TOKEN")

	return nil
}

type GuildConfig struct {
	AfkChannelId           string
	GuildId                string
	Prefix                 string
	PlaylistChannelId      string
	PlaylistMessageImageId string
	PlaylistMessageQueueId string
}
