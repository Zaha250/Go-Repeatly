package main

import (
	"context"
	"log"
	"repeatly/internal/app"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Ошибка при инициализации приложения: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
