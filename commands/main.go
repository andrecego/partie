package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isCommand(message string) bool {
	return strings.HasPrefix(message, "!")
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
