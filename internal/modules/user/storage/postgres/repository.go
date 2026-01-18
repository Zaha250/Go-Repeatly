package postgres

import (
	"context"
	"database/sql"
	"errors"
	"repeatly/internal/database"
	"repeatly/internal/modules/user/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db}
}

func (repo *UserRepository) FindByTelegramID(ctx context.Context, tgID int64) (*domain.User, error) {
	query := `SELECT id, telegram_id, first_name, username, status, timezone, created_at FROM users WHERE telegram_id = $1`

	var user domain.User
	err := repo.db.QueryRowContext(ctx, query, tgID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.FirstName,
		&user.Username,
		&user.Status,
		&user.Timezone,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (telegram_id, first_name, username, timezone) VALUES ($1, $2, $3, $4)`
	_, err := repo.db.ExecContext(ctx, query, user.TelegramID, user.FirstName, user.Username, user.Timezone)
	return err
}
