package telegram

import (
	"context"
	"fmt"
	"log"
	"repeatly/internal/pkg/telegram"

	userDomain "repeatly/internal/modules/user/domain"
	userUsecase "repeatly/internal/modules/user/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Handler обрабатывает входящие обновления от Telegram.
type Handler struct {
	userService *userUsecase.UserService
	// ... сюда будем добавлять другие сервисы (noteService и т.д.)
}

const (
	WelcomeMessage     = "Добро пожаловать, %s! 👋\n\nЯ помогу вам запоминать любую информацию с помощью кривой забывания Эббингауза. Просто отправьте мне текст, фото или файл, который хотите выучить."
	WelcomeBackMessage = "С возвращением, %s! 💪\n\nГотовы продолжить? Отправьте мне новую информацию для запоминания."
)

func NewHandler(userService *userUsecase.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) HandleUpdate(ctx context.Context, bot telegramClient.BotClient, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		h.handleCommand(ctx, bot, update.Message)
		return
	}

	log.Printf("Received a non-command message from %s", update.Message.From.UserName)
}

func (h *Handler) handleCommand(ctx context.Context, bot telegramClient.BotClient, message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		h.handleStartCommand(ctx, bot, message)
	default:
		// Отправляем сообщение о неизвестной команде
		if err := bot.SendMessage(message.Chat.ID, "Я не знаю такой команды."); err != nil {
			log.Printf("Error sending 'unknown command' message: %v", err)
		}
	}
}

func (h *Handler) handleStartCommand(ctx context.Context, bot telegramClient.BotClient, message *tgbotapi.Message) {
	user := &userDomain.User{
		TelegramID: message.From.ID,
		FirstName:  message.From.FirstName,
		Username:   message.From.UserName,
		Timezone:   "+00", //временный хардкод
	}

	_, wasCreated, err := h.userService.RegisterIfNotExist(ctx, user)
	if err != nil {
		log.Printf("Error during user registration: %v", err)
		bot.SendMessage(message.Chat.ID, "Ой, что-то пошло не так. Попробуйте позже.")
		return
	}

	var msgText string
	if wasCreated {
		msgText = fmt.Sprintf(WelcomeMessage, user.FirstName)
	} else {
		msgText = fmt.Sprintf(WelcomeBackMessage, user.FirstName)
	}

	if err := bot.SendMessage(message.Chat.ID, msgText); err != nil {
		log.Printf("Error sending welcome message: %v", err)
	}
}
