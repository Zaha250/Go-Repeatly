package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"repeatly/internal/database"
	"repeatly/internal/delivery/telegram"
	"repeatly/internal/pkg/telegram"
	"repeatly/internal/server"
	"syscall"

	"github.com/joho/godotenv"

	userHTTP "repeatly/internal/modules/user/delivery/http"
	userPostgres "repeatly/internal/modules/user/storage/postgres"
	userUsecase "repeatly/internal/modules/user/usecase"
	"repeatly/internal/pkg/config"
)

// App инкапсулирует все компоненты приложения.
type App struct {
	cfg        *config.Config
	db         *database.DB
	bot        telegramClient.BotClient
	botHandler *telegram.Handler
	httpServer *http.Server
}

// NewApp создает и инициализирует новый экземпляр приложения.
func NewApp(ctx context.Context) (*App, error) {
	// 1. Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 2. Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Не удалось загрузить файл конфигурации: %v", err)
	}

	// 3. Подключение к БД
	db, err := database.ConnectDB(&cfg.Postgres)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	log.Println("Успешное подключение к базе данных")

	// 4. Инициализация Telegram Client Bot
	bot, err := telegramClient.NewClient(&cfg.Telegram)
	if err != nil {
		return nil, fmt.Errorf("ошибка при инициализации telegram bot: %w", err)
	}

	// 5. --- Сборочный корень (Dependency Injection) ---
	userRepo := userPostgres.NewUserRepository(db)
	userService := userUsecase.NewUserService(userRepo)

	// 6. Создание обработчиков
	botHandler := telegram.NewHandler(userService)
	userHandler := userHTTP.NewUserHandler(userService)

	routerFactory := server.NewRouterFactory(userHandler)
	router := routerFactory.InitRouter()

	httpServer := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	return &App{
		cfg:        cfg,
		db:         db,
		bot:        bot,
		botHandler: botHandler,
		httpServer: httpServer,
	}, nil
}

// Run запускает все долгоживущие процессы (слушатель бота, http-сервер).
func (a *App) Run(ctx context.Context) error {
	// Создаем контекст для отслеживания сигналов завершения (Ctrl+C)
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("Сервер запущен на порту %s", a.cfg.App.Port)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске HTTP сервера: %v", err)
		}
	}()

	// Запускаем слушателя Telegram в отдельной горутине
	go a.runTelegramListener(ctx)

	// Ждем сигнала завершения
	<-ctx.Done()
	log.Println("Shutdown signal received")

	// Выполняем graceful shutdown
	return a.shutdown()
}

func (a *App) runTelegramListener(ctx context.Context) {
	updates, err := a.bot.GetUpdatesChan(60) // <-- ИЗМЕНЕНИЕ
	if err != nil {
		log.Printf("Failed to get updates channel: %v", err)
		return
	}

	log.Println("Telegram listener started")

	for {
		select {
		case <-ctx.Done(): // Если пришел сигнал завершения, выходим
			log.Println("Telegram listener stopped")
			return
		case update := <-updates:
			// Передаем обработку в наш новый хендлер
			a.botHandler.HandleUpdate(ctx, a.bot, update)
		}
	}
}

// shutdown корректно завершает работу всех компонентов.
func (a *App) shutdown() error {
	log.Println("Shutting down application...")
	if err := a.db.Close(); err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}
	log.Println("Database connection closed.")
	return nil
}
