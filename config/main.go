package config

import "os"

var (
	Token     string //To store value of Token from config.json .
	BotPrefix string // To store value of BotPrefix from config.json.
)

func ReadConfig() error {
	Token = os.Getenv("TOKEN")
	BotPrefix = "!"

	return nil
}
