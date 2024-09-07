package helpers

import (
	"partie-bot/config" //importing our config package which we have created above
	"strings"
)

func MsgWithPrefix(name string, guildConfig config.GuildConfig) string {
	return guildConfig.Prefix + name
}

func ParseCommand(content string, config config.GuildConfig) (string, []string) {
	if !strings.HasPrefix(content, config.Prefix) {
		return "", nil
	}

	content = strings.TrimPrefix(content, config.Prefix)
	words := strings.Fields(content)
	if len(words) == 0 {
		return "", nil
	}

	return words[0], words[1:]
}
