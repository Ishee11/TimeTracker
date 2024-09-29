package router

import (
	"time-tracker/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршрутизатор и возвращает готовый http.Handler
func SetupRouter() *gin.Engine {
	r := gin.Default() // Создаем новый маршрутизатор с готовыми middleware

	// Настраиваем обработчик для добавления активности
	r.POST("/activities", handlers.AddActivityHandler)

	// Вы можете добавлять другие маршруты таким же образом
	// r.POST("/another-endpoint", someOtherHandlerFunc)

	return r
}
