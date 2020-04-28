package habrbestbot

import (
	"context"
	"log"
	"os"

	"github.com/baor/habr-best-bot/habr"
	"github.com/baor/habr-best-bot/storage"
	"github.com/baor/habr-best-bot/telegram"
)

type botContext struct {
	tlg        telegram.Messenger
	tlgChannel string
	st         storage.PostStorer
	feed       habr.FeedReader
}

// full list of allowed tags []string{"a", "b", "strong", "i", "em", "code", "pre"}
// I will leave only a to avoid tag nesting
var telegramAllowedTags = []string{"a"}

func (c *botContext) updateFeedToChannel() {
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

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// HelloPubSub consumes a Pub/Sub message.
func Entrypoint(ctx context.Context, m PubSubMessage) error {
	log.Println(string(m.Data))

	token := os.Getenv("TELEGRAM_API_TOKEN")
	if token == "" {
		log.Panic("Empty TELEGRAM_API_TOKEN")
	}
	bot := telegram.NewBot(token)

	const GCS_BUCKET_NAME = "habr-best-feeds-storage-2"
	log.Printf("GCS bucket name: %s", GCS_BUCKET_NAME)

	s := storage.NewGcsAdapter(GCS_BUCKET_NAME)

	bCtx := botContext{
		tlg:        bot,
		tlgChannel: "@habrbest",
		st:         s,
		feed:       habr.NewHabrReader(),
	}

	bCtx.updateFeedToChannel()
	log.Println("Feed was updated successfully")
	return nil
}
