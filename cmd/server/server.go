package server

import (
	"log"
	"net/http"
	"repeatly/internal/pkg/config"

	"repeatly/internal/database"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg *config.Config
	db  *database.DB
}

func NewServer(cfg *config.Config, db *database.DB) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
	}
}

func (s *Server) setupRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}

func (s *Server) Run() error {
	router := s.setupRoutes()
	port := s.cfg.App.Port
	log.Printf("Сервер запущен на порту %s", port)
	return router.Run(":" + port)
}
