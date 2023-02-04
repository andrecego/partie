package config

import "os"

var (
	Token        string // To store value of Token from config.json .
	BotPrefix    string // To store value of BotPrefix from config.json.
	RollbarToken string
)

func ReadConfig() error {
	Token = os.Getenv("TOKEN")
	RollbarToken = os.Getenv("ROLLBAR_TOKEN")
	BotPrefix = "!"

	return nil
}
