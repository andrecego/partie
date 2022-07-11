package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var prefixlessCommands = []string{
	"skip",
	"pause",
	"play",
	"restart",
	"remove",
	"delete",
}

func isCommand(message string) bool {
	return strings.HasPrefix(message, "!")
}

func isPrefixlessCommands(message string) bool {
	for _, command := range prefixlessCommands {
		if strings.HasPrefix(message, command) {
			return true
		}
	}

	return false
}

func commandParse(message string) (string, []string) {
	args := strings.Split(message, " ")
	return args[0], args[1:]
}

func parseCommand(command string) string {
	return command[1:]
}

func tagUser(user *discordgo.User) string {
	return "<@" + user.ID + ">"
}

func parseCommands(m *discordgo.MessageCreate, command string) (string, []string) {
	if !isCommand(m.Content) {
		return "", nil
	}

	commands := m.Content[1:]
	args := strings.Split(commands, " ")

	if args[0] != command {
		return "", nil
	}

	return args[0], args[1:]
}
