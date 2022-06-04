package music

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/bwmarrin/discordgo"
)

const maxQueuePrintSize = 20

func nowPlayingMessage() {
	channelID := "943655307626823771"
	song := currentDJ.CurrentSong

	message := fmt.Sprintf("Now playing: %s\nAdded by: %s", song.GetTitle(), song.GetAddedBy())
	currentDJ.Discord.Session.ChannelMessageSend(channelID, message)
}

func getMessage() (string, *discordgo.MessageEmbed) {
	npMessage := ""
	npStatus := ""
	npImage := "https://media.wired.com/photos/5ee79f9d6f93e879afd83412/191:100/w_2400,h_1256,c_limit/Science_Spot.jpg"
	np := currentDJ.CurrentSong
	if np != nil {
		npImage = np.GetThumbnail()
		npStatus = "Now playing:"
		npMessage = embedFormatSong(np)
	} else {
		npStatus = "Nothing playing."
		npMessage = "To play a song just type the name of the song"
	}

	queueMessage := "__**Queue:**__\n"
	if len(currentDJ.Queue) > 0 {
		queueSize := len(currentDJ.Queue)
		var printedQueue []Song

		if queueSize > maxQueuePrintSize {
			queueMessage += fmt.Sprintf("And %d more songs\n", queueSize-maxQueuePrintSize)
			queueSize = maxQueuePrintSize

			printedQueue = make([]Song, maxQueuePrintSize)
			copy(printedQueue, currentDJ.Queue[:maxQueuePrintSize])
		} else {
			printedQueue = make([]Song, queueSize)
			copy(printedQueue, currentDJ.Queue)
		}
		ReverseSlice(printedQueue)
		for i, song := range printedQueue {
			queueMessage += fmt.Sprintf("%02d. %s\n", queueSize-i, formatSong(song))
		}

	} else {
		queueMessage += "No songs in queue. Go add some! 🎵"
	}

	embededMessage := &discordgo.MessageEmbed{Title: npStatus, Description: npMessage, Color: 0x5E06AC}
	embededMessage.Image = &discordgo.MessageEmbedImage{URL: npImage}
	embededMessage.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("!music commands work here as well"), IconURL: "https://www.pngmart.com/files/11/Doge-Meme-PNG-Photos.png"}
	return queueMessage, embededMessage
}

func queueMessage(channelID string) {
	queueMessage, embededMessage := getMessage()
	message := &discordgo.MessageSend{Content: queueMessage, Embed: embededMessage, AllowedMentions: nil}

	currentDJ.Discord.Session.ChannelMessageSendComplex(channelID, message)
}

func addedToQueueMessage(song Song) {
	channelID := song.GetChannelID()
	if channelID == "955146633203560468" { // playlist channel id
		return
	}

	message := fmt.Sprintf("Added %s to queue.\n", song.GetTitle())
	currentDJ.Discord.Session.ChannelMessageSend(channelID, message)
}

func formatSong(song Song) string {
	return fmt.Sprintf("%s [%s] - %s", truncate(song.GetTitle(), 50), fmtDuration(song.GetDuration()), song.GetAddedBy())
}

func embedFormatSong(song Song) string {
	return fmt.Sprintf("[%s](%s) [%s] - %s", truncate(song.GetTitle(), 50), song.GetVideoURL(), fmtDuration(song.GetDuration()), song.GetAddedBy())
}

func fmtDuration(d time.Duration) string {
	seconds := d.Round(time.Second)
	minutes := seconds / time.Minute
	d -= minutes * time.Minute
	seconds = d / time.Second
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func truncate(s string, maxLen int) string {
	length := len(s)
	if length <= maxLen {
		return s
	}

	truncated := ""
	count := 0
	for _, char := range s {
		truncated += string(char)
		if count >= length {
			break
		}
		count++
	}
	return truncated
}

func numberOfSongsInQueue() string {
	size := len(currentDJ.Queue)
	if size == 0 {
		return "No"
	}

	return fmt.Sprintf("%02d", size)
}

func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func updateQueueMessage() {
	messageID := "955235598292103199"
	channelID := "955146633203560468"

	content, embed := getMessage()

	messageEdit := &discordgo.MessageEdit{
		ID:              messageID,
		Channel:         channelID,
		Content:         &content,
		Embed:           embed,
		AllowedMentions: nil,
	}

	_, err := currentDJ.Discord.Session.ChannelMessageEditComplex(messageEdit)
	if err != nil {
		fmt.Println("Error editing message: ", err)
	}
}
