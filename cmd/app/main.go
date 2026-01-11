package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"repeatly/internal/pkg/config"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // Анонимный импорт драйвера
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить файл конфигурации: %v", err)
	}

	db, err := connectDB(&cfg.Postgres)
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

func connectDB(cfg *config.PostgresConfig) (*sql.DB, error) {
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

	return db, nil
}
