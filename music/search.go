package music

import (
	"fmt"
	"partie-bot/repositories"

	"github.com/bwmarrin/discordgo"
)

func Search(query string, message *discordgo.MessageCreate) {
	results, err := repositories.YoutubeSearch(query)
	if err != nil {
		fmt.Println("Error searching youtube: ", err)
		return
	}

	if len(results) == 0 {
		currentDJ.Discord.Session.ChannelMessageSend(message.ChannelID, "No results found.")
		return
	}

	var filterResults []repositories.YoutubeSearchResult
	for _, result := range results {
		if result.Duration() == 0 || result.Duration() > 600 {
			continue
		}

		filterResults = append(filterResults, result)
	}

	respMessage := "Results:\n"
	for i, result := range filterResults {
		respMessage += fmt.Sprintf("%d. (%s) %s\n", i, result.DurationRaw, result.Title)
	}

	err = repositories.Save(message.GuildID, message.Author.ID, filterResults)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentDJ.Discord.Session.ChannelMessageSend(message.ChannelID, respMessage)
}
