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

func (c *botContext) updateFeedToChannel() {
	for _, feedItem := range c.feed.GetBestFeed() {
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

// pubSubMessage is the payload of a Pub/Sub event.
type pubSubMessage struct {
	Data []byte `json:"data"`
}

// Entrypoint consumes a Pub/Sub message which triggers feed update.
func Entrypoint(ctx context.Context, m pubSubMessage) error {
	log.Println(string(m.Data))

	token := os.Getenv("TELEGRAM_API_TOKEN")
	if token == "" {
		log.Panic("Empty TELEGRAM_API_TOKEN")
	}
	bot := telegram.NewBot(token)

	GCS_BUCKET_NAME := os.Getenv("GCS_BUCKET_NAME")
	if GCS_BUCKET_NAME == "" {
		log.Panic("Empty GCS_BUCKET_NAME")
	}
	log.Printf("GCS bucket name: %s", GCS_BUCKET_NAME)
	s := storage.NewGcsAdapter(GCS_BUCKET_NAME)

	// FIRESTORE_CLOUD_PROJECT := os.Getenv("FIRESTORE_CLOUD_PROJECT")
	// if FIRESTORE_CLOUD_PROJECT == "" {
	// 	log.Panic("Empty FIRESTORE_CLOUD_PROJECT")
	// }
	// FIRESTORE_COLLECTION_NAME := os.Getenv("FIRESTORE_COLLECTION_NAME")
	// if FIRESTORE_COLLECTION_NAME == "" {
	// 	log.Panic("Empty FIRESTORE_COLLECTION_NAME")
	// }
	// s := storage.NewFirestoreAdapter(FIRESTORE_COLLECTION_NAME, FIRESTORE_CLOUD_PROJECT)

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
