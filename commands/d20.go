package commands

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var d20Resuls = map[int][]string{
	1:  {"You rolled a 1. You are a loser."},
	2:  {"You rolled a 2. You are a loser."},
	3:  {"You rolled a 3. You are a loser."},
	4:  {"You rolled a 4. You are a loser."},
	5:  {"You rolled a 5. You are a loser."},
	6:  {"You rolled a 6. You are a loser."},
	7:  {"You rolled a 7. You are a loser."},
	8:  {"You rolled a 8. You are a loser."},
	9:  {"You rolled a 9. You are a loser."},
	10: {"You rolled a 10. You are a loser."},
	11: {"You rolled a 11. You are a winner."},
	12: {"You rolled a 12. You are a winner."},
	13: {"You rolled a 13. You are a winner."},
	14: {"You rolled a 14. You are a winner."},
	15: {"You rolled a 15. You are a winner."},
	16: {"You rolled a 16. You are a winner."},
	17: {"You rolled a 17. You are a winner."},
	18: {"You rolled a 18. You are a winner."},
	19: {"You rolled a 19. You are a winner."},
	20: {"You rolled a 20. You are a winner."},
}

// var d20Commands = map[int]Handler{}

// type Handler interface {
// 	Execute(s *discordgo.Session, m *discordgo.MessageCreate)
// }

func RollD20Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isCommand(m.Content) {
		return
	}

	command := parseCommand(m.Content)
	if command != "d20" {
		return
	}

	rand.Seed(time.Now().UnixNano())

	roll := rand.Intn(20) + 1

	s.ChannelMessageSend(m.ChannelID, tagUser(m.Author)+" "+d20Resuls[roll][0])
}
