package telegram

type Messenger interface {
	NewMessageToChannel(username string, text string) error
}
