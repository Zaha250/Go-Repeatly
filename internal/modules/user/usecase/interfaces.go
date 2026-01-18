package usecase

import (
	"context"
	"repeatly/internal/modules/user/domain"
)

type UserRepository interface {
	// FindByTelegramID ищет пользователя по его ID в Telegram.
	FindByTelegramID(ctx context.Context, tgID int64) (*domain.User, error)
	// Create создает нового пользователя в хранилище.
	Create(ctx context.Context, user *domain.User) error
}
