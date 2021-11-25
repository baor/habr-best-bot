package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
}

// NewBot returns an instance of Bot which implements Messenger interface
func NewBot(token string) Messenger {
	var err error
	b := Bot{}
	b.botAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	b.botAPI.Debug = false
	log.Printf("Telegram bot authorized on account %s", b.botAPI.Self.UserName)
	return &b
}

func (b *Bot) NewMessageToChannel(username string, text string) error {
	MAX_LENGTH := 4090
	if len(text) > MAX_LENGTH {
		log.Printf("trim too long string %s", text)
		text = text[:MAX_LENGTH] + "..."
	}
	msg := tgbotapi.NewMessageToChannel(username, text)
	msg.ParseMode = "HTML"
	_, err := b.botAPI.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
