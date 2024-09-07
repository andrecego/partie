package helpers

import (
	"partie-bot/config"
	"testing"
)

// Test for ParseCommand function
func TestParseCommand(t *testing.T) {
	tests := []struct {
		scenario     string
		content      string
		guildConfig  config.GuildConfig
		expectedCmd  string
		expectedArgs []string
	}{
		{
			scenario:     "Content does not start with prefix",
			content:      "command args",
			guildConfig:  config.GuildConfig{Prefix: "!"},
			expectedCmd:  "",
			expectedArgs: nil,
		},
		{
			scenario:     "Content starts with prefix but empty after prefix",
			content:      "!",
			guildConfig:  config.GuildConfig{Prefix: "!"},
			expectedCmd:  "",
			expectedArgs: nil,
		},
		{
			scenario:     "Content starts with prefix and has one word",
			content:      "!command",
			guildConfig:  config.GuildConfig{Prefix: "!"},
			expectedCmd:  "command",
			expectedArgs: []string{},
		},
		{
			scenario:     "Content starts with prefix and has multiple words",
			content:      "!command arg1 arg2",
			guildConfig:  config.GuildConfig{Prefix: "!"},
			expectedCmd:  "command",
			expectedArgs: []string{"arg1", "arg2"},
		},
		{
			scenario:     "Prefix is not a single character",
			content:      "!!!command arg1 arg2",
			guildConfig:  config.GuildConfig{Prefix: "!!"},
			expectedCmd:  "!command",
			expectedArgs: []string{"arg1", "arg2"},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			cmd, args := ParseCommand(test.content, test.guildConfig)
			if cmd != test.expectedCmd {
				t.Errorf("Expected command %s but got %s", test.expectedCmd, cmd)
			}
			if !equalSlices(args, test.expectedArgs) {
				t.Errorf("Expected args %v but got %v", test.expectedArgs, args)
			}
		})
	}
}

// Helper function to compare string slices
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
