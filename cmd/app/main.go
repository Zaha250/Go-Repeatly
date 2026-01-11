package main

import (
	"log"
	"net/http"
	"repeatly/internal/database"
	"repeatly/internal/pkg/config"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	// Тестовый эндпоинт, чтобы проверить, что сервер работает
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Printf("Сервер запущен на порту %s", cfg.App.Port)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
