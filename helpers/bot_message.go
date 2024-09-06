package helpers

import (
	"partie-bot/config" //importing our config package which we have created above
)

func MsgWithPrefix(name string) string {
	return config.BotPrefix + name
}
