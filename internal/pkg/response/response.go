package response

import "github.com/gin-gonic/gin"

// SuccessResponse определяет структуру успешного ответа.
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse определяет структуру ответа с ошибкой.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success - это универсальная функция для отправки успешного ответа.
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		Data: data,
	})
}

// Error - это универсальная функция для отправки ответа с ошибкой.
func Error(c *gin.Context, statusCode int, message string) {
	// Здесь можно добавить логирование ошибки
	c.JSON(statusCode, ErrorResponse{
		Error: message,
	})
}
