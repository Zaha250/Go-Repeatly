package main

import (
	"log"
	"repeatly/cmd/server"
	"repeatly/internal/database"
	"repeatly/internal/pkg/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить файл конфигурации: %v", err)
	}

	db, err := database.ConnectDB(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	log.Println("Успешное подключение к базе данных")

	srv := server.NewServer(cfg, db)
	if err := srv.Run(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
