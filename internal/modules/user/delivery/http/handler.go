package http

import (
	"errors"
	"net/http"
	"repeatly/internal/modules/user/storage/postgres"
	userUsecase "repeatly/internal/modules/user/usecase"
	"repeatly/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *userUsecase.UserService
}

func NewUserHandler(service *userUsecase.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		users := api.Group("/user")
		{
			users.GET("/by-telegram-id/:id", h.getUserByTelegramID)
		}
	}
}

func (h *UserHandler) getUserByTelegramID(c *gin.Context) {
	idStr := c.Param("id")
	tgID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid Telegram ID format")
		return
	}

	user, err := h.userService.FindByTelegramID(c.Request.Context(), tgID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "Пользователь не найден")
			return
		}
		// Для всех других ошибок - 500
		response.Error(c, http.StatusInternalServerError, "Internal server error")
		return
	}
	response.Success(c, http.StatusOK, user)
}
