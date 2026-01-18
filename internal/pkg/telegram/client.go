package telegramClient

import (
	"repeatly/internal/pkg/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotClient interface {
	SendMessage(chatID int64, text string) error
	GetUpdatesChan(timeout int) (tgbotapi.UpdatesChannel, error)
}

type apiWrapper struct {
	bot *tgbotapi.BotAPI
}

// NewClient создает новый экземпляр нашего клиента.
func NewClient(telegramConfig *config.TelegramConfig) (BotClient, error) {
	bot, err := tgbotapi.NewBotAPI(telegramConfig.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = telegramConfig.Debug

	return &apiWrapper{bot: bot}, nil
}

// SendMessage реализует метод интерфейса, вызывая под капотом метод библиотеки.
func (w *apiWrapper) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := w.bot.Send(msg)
	return err
}

// GetUpdatesChan реализует метод интерфейса.
func (w *apiWrapper) GetUpdatesChan(timeout int) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout
	return w.bot.GetUpdatesChan(u), nil
}
