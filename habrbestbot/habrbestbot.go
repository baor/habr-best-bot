package habrbestbot

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/baor/habr-best-bot/habr"
	"github.com/baor/habr-best-bot/storage"
	"github.com/baor/habr-best-bot/telegram"
)

type context struct {
	tlg        telegram.Messenger
	tlgChannel string
	st         storage.PostStorer
	feed       habr.FeedReader
}

// full list of allowed tags []string{"a", "b", "strong", "i", "em", "code", "pre"}
// I will leave only a to avoid tag nesting
var telegramAllowedTags = []string{"a"}

func (c *context) updateFeedToChannel() {
	for _, feedItem := range c.feed.GetBestFeed(telegramAllowedTags) {
		if c.st.IsPostIDExists(feedItem.ID) {
			continue
		}
		log.Printf("Send post %s", feedItem.ID)
		err := c.tlg.NewMessageToChannel(c.tlgChannel, feedItem.Message)
		if err != nil {
			log.Panic(err)
		}
		c.st.AddPostID(feedItem.ID)
	}
}

func entrypoint(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("TELEGRAM_API_TOKEN")
	if token == "" {
		log.Panic("Empty TELEGRAM_API_TOKEN")
	}
	bot := telegram.NewBot(token)

	gcsBucketName := "habrfeeds"
	log.Printf("GCS bucket name: %s", gcsBucketName)
	s := storage.NewGcsAdapter(gcsBucketName)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	ctx := context{
		tlg:        bot,
		tlgChannel: "@habrbest",
		st:         s,
		feed:       habr.NewHabrReader(),
	}

	ctx.updateFeedToChannel()
	fmt.Fprint(w, "Feed was updated successfully")
}
