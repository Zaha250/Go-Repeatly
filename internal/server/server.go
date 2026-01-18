package server

import (
	"net/http"
	userHTTP "repeatly/internal/modules/user/delivery/http" // Импортируем наш HTTP-хендлер

	"github.com/gin-gonic/gin"
)

// RouterFactory инкапсулирует создание роутера и регистрацию маршрутов.
type RouterFactory struct {
	userHandler *userHTTP.UserHandler
}

func NewRouterFactory(userHandler *userHTTP.UserHandler) *RouterFactory {
	return &RouterFactory{
		userHandler: userHandler,
	}
}

func (f *RouterFactory) InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	f.userHandler.RegisterRoutes(router)

	return router
}
