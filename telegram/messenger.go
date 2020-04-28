package telegram

type Messenger interface {
	NewMessageToChat(chatID int64, text string) error
	NewMessageToChannel(username string, text string) error
}
