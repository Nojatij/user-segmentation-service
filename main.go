// main.go
package main

import (
	"user-segmentation-service/api"
	"user-segmentation-service/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db.InitDB()

	// Создание нового маршрутизатора Gin
	r := gin.Default()

	// Настройка маршрутов
	r.POST("/segments", api.CreateSegment)
	r.DELETE("/segments/:slug", api.DeleteSegment)
	r.POST("/segments/user", api.AddUserToSegment)
	r.GET("/segments/user", api.GetActiveSegments)

	// Запуск сервера
	r.Run(":8080")
}
