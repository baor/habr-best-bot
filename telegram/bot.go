package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type updateMessage struct {
	Text   string
	ChatID int64
}

type updatesChannel <-chan updateMessage

type Bot struct {
	botAPI  *tgbotapi.BotAPI
	updates updatesChannel
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

func (b *Bot) NewMessageToChat(chatID int64, text string) error {
	if len(text) > 4096 {
		log.Printf("trim too long string %s", text)
		text = text[:4090] + "..."
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	log.Println(msg)
	_, err := b.botAPI.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) NewMessageToChannel(username string, text string) error {
	if len(text) > 4096 {
		log.Printf("trim too long string %s", text)
		text = text[:4090] + "..."
	}
	msg := tgbotapi.NewMessageToChannel(username, text)
	msg.ParseMode = "HTML"
	log.Println(msg)
	_, err := b.botAPI.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
