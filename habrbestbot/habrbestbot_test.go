package habrbestbot

import (
	"testing"

	"github.com/baor/habr-best-bot/habr"
	"github.com/stretchr/testify/assert"
)

func TestHabrBestBot_updateFeedToChannel_OneItem(t *testing.T) {
	bot := telegramBotMocked{}
	s := storageMocked{internalStorage: map[string]bool{}}
	fr := feedReaderMocked{
		items: []habr.FeedItem{
			{
				LinkToImage: "http://",
				Message:     "text",
				ID:          "id1",
			},
		},
	}
	ctx := context{
		tlg:        &bot,
		tlgChannel: "@habrbest",
		st:         &s,
		feed:       fr,
	}

	ctx.updateFeedToChannel()
	assert.Equal(t, 1, bot.newMessagesCount)
	assert.True(t, s.internalStorage["id1"])
}

func TestHabrBestBot_updateFeedToChannel_TwoItems(t *testing.T) {
	bot := telegramBotMocked{}
	s := storageMocked{internalStorage: map[string]bool{}}
	fr := feedReaderMocked{
		items: []habr.FeedItem{
			{
				LinkToImage: "http://",
				Message:     "text",
				ID:          "id1",
			},
			{
				LinkToImage: "http://",
				Message:     "text",
				ID:          "id2",
			},
		},
	}
	ctx := context{
		tlg:        &bot,
		tlgChannel: "@habrbest",
		st:         &s,
		feed:       fr,
	}

	ctx.updateFeedToChannel()
	assert.Equal(t, 2, bot.newMessagesCount)
	assert.True(t, s.internalStorage["id1"])
	assert.True(t, s.internalStorage["id2"])
}

func TestHabrBestBot_updateFeedToChannel_TwoSameItems(t *testing.T) {
	bot := telegramBotMocked{}
	s := storageMocked{internalStorage: map[string]bool{}}
	fr := feedReaderMocked{
		items: []habr.FeedItem{
			{
				LinkToImage: "http://",
				Message:     "text",
				ID:          "id1",
			},
			{
				LinkToImage: "http://",
				Message:     "text",
				ID:          "id1",
			},
		},
	}
	ctx := context{
		tlg:        &bot,
		tlgChannel: "@habrbest",
		st:         &s,
		feed:       fr,
	}

	ctx.updateFeedToChannel()
	assert.Equal(t, 1, bot.newMessagesCount)
	assert.True(t, s.internalStorage["id1"])
}

type telegramBotMocked struct {
	newMessagesCount int
}

func (t *telegramBotMocked) NewMessageToChat(chatID int64, text string) error {
	t.newMessagesCount++
	return nil
}

func (t *telegramBotMocked) NewMessageToChannel(username string, text string) error {
	t.newMessagesCount++
	return nil
}

type storageMocked struct {
	internalStorage map[string]bool
}

func (s storageMocked) IsPostIDExists(id string) bool {
	return s.internalStorage[id]
}
func (s storageMocked) AddPostID(id string) {
	s.internalStorage[id] = true
}

type feedReaderMocked struct {
	items []habr.FeedItem
}

func (fr feedReaderMocked) GetBestFeed(allowedTags []string) []habr.FeedItem {
	return fr.items
}
