package database

import (
	"database/sql"
	"fmt"
	"repeatly/internal/pkg/config"
)

type DB struct {
	*sql.DB
}

func ConnectDB(cfg *config.PostgresConfig) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db}, nil
}
