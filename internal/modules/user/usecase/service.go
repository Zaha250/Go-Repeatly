package usecase

import (
	"context"
	"errors"
	"log"
	"repeatly/internal/modules/user/domain"
	"repeatly/internal/modules/user/storage/postgres"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo}
}

func (s *UserService) RegisterIfNotExist(ctx context.Context, tgUser *domain.User) (*domain.User, bool, error) {
	//1. Пытаемся найти пользователя
	existingUser, err := s.userRepo.FindByTelegramID(ctx, tgUser.TelegramID)
	if err != nil && !errors.Is(err, postgres.ErrUserNotFound) {
		log.Printf("Error finding user: %v", err)
		return nil, false, err
	}

	//2. Если находим, то возвращаем его
	if existingUser != nil {
		return existingUser, false, nil
	}

	//3. Если НЕ находим, регистрируем нового
	log.Printf("Пользователь с telegram_id %d не найден, регистрируем", tgUser.TelegramID)
	newUser := &domain.User{
		TelegramID: tgUser.TelegramID,
		FirstName:  tgUser.FirstName,
		Username:   tgUser.Username,
		Timezone:   tgUser.Timezone,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		log.Printf("Ошибка регистрации пользователя: %v", err)
		return nil, false, err
	}

	createdUser, err := s.userRepo.FindByTelegramID(ctx, tgUser.TelegramID)
	if err != nil {
		return nil, true, err
	}
	return createdUser, true, nil
}

func (s *UserService) FindByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	user, err := s.userRepo.FindByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
