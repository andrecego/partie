package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPrefilessCommand(t *testing.T) {
	assert.True(t, isPrefixlessCommands("skip"))
	assert.True(t, isPrefixlessCommands("play"))
	assert.True(t, isPrefixlessCommands("pause"))
	assert.True(t, isPrefixlessCommands("remove 1"))
	assert.False(t, isPrefixlessCommands(" skip"))
	assert.False(t, isPrefixlessCommands("!skip"))
}
