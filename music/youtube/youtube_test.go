package youtube

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func loadFixture(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func testYoutube_setup(t *testing.T) Youtube {
	data, err := os.ReadFile("../../fixtures/test_song.json")
	assert.NoError(t, err)

	var result YoutubeResult
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)
	return result.Entries[0]
}

func TestYoutube_GetID(t *testing.T) {
	yt := testYoutube_setup(t)

	assert.Equal(t, "7XVG4oWQbx4", yt.GetID())
}

func TestYoutube_GetTitle(t *testing.T) {
	yt := testYoutube_setup(t)

	assert.Equal(t, "Bruno Mars - That’s What I Like ", yt.GetTitle())
}

func TestYoutube_GetDuration(t *testing.T) {
	yt := testYoutube_setup(t)

	assert.Equal(t, 238*time.Second, yt.GetDuration())
}

func TestYoutube_GetAddedBy(t *testing.T) {
	yt := testYoutube_setup(t)
	user := &discordgo.User{ID: "user-id"}
	yt.AddedBy = AddedBy{User: user}
	assert.Equal(t, "<@user-id>", yt.GetAddedBy())
}

func TestYoutube_GetAuthorID(t *testing.T) {
	yt := testYoutube_setup(t)
	user := &discordgo.User{ID: "user-id"}
	yt.AddedBy = AddedBy{User: user}
	assert.Equal(t, "user-id", yt.GetAuthorID())
}

func TestYoutube_GetChannelID(t *testing.T) {
	yt := testYoutube_setup(t)
	channel := &discordgo.Channel{ID: "channel-id"}
	yt.AddedBy = AddedBy{Channel: channel}
	assert.Equal(t, "channel-id", yt.GetChannelID())
}

func TestYoutube_GetGuildID(t *testing.T) {
	yt := testYoutube_setup(t)
	guild := &discordgo.Guild{ID: "guild-id"}
	yt.AddedBy = AddedBy{Guild: guild}
	assert.Equal(t, "guild-id", yt.GetGuildID())
}

func TestYoutube_GetURL(t *testing.T) {
	yt := testYoutube_setup(t)
	assert.Equal(t, "https://rr6---sn-8p8v-bg0lr.googlevideo.com/videoplayback?expire=1727631227&ei=Gzv5Zr6iHJCOobIP7rrkkAY&ip=152.249.28.247&id=o-ALCYhJ1OOm5Bj-ocdSL9HID3WTtPCGjALzJ_n0lfrD5h&itag=251&source=youtube&requiressl=yes&xpc=EgVo2aDSNQ%3D%3D&mh=-O&mm=31%2C29&mn=sn-8p8v-bg0lr%2Csn-8p8v-bg0s6&ms=au%2Crdu&mv=m&mvi=6&pcm2cms=yes&pl=24&initcwndbps=1311250&vprv=1&svpuc=1&mime=audio%2Fwebm&rqh=1&gir=yes&clen=3836160&dur=237.781&lmt=1627183832407342&mt=1727609118&fvip=2&keepalive=yes&fexp=51299152&c=IOS&txp=5432434&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cxpc%2Cvprv%2Csvpuc%2Cmime%2Crqh%2Cgir%2Cclen%2Cdur%2Clmt&sig=AJfQdSswRgIhAOONh-Qyj6MDYiJ3fhGpRePLUaCu1ci6fpvqRb_4WYxXAiEAj-5KZeDaAl1357SoXOOj4kbaENFjRrbGILIKCqUQcOc%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpcm2cms%2Cpl%2Cinitcwndbps&lsig=ABPmVW0wRAIgIPmbL6ChApBiezZ1o4m71HDrzMQ58nDZ7zlNyu4dFZQCID0A7d6p501MB29kamUeVznxXcAsOyX1TKkCg8LLA1b1", yt.GetURL())
}

func TestYoutube_GetThumbnail(t *testing.T) {
	yt := testYoutube_setup(t)
	assert.Equal(t, "https://i.ytimg.com/vi_webp/7XVG4oWQbx4/maxresdefault.webp", yt.GetThumbnail())
}

func TestYoutube_GetStartTime(t *testing.T) {
	yt := testYoutube_setup(t)
	assert.Equal(t, 0, yt.GetStartTime())
}

func TestYoutube_GetVideoURL(t *testing.T) {
	yt := testYoutube_setup(t)
	assert.Equal(t, "https://www.youtube.com/watch?v=7XVG4oWQbx4", yt.GetVideoURL())
}
