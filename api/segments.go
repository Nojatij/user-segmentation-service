// api/segments.go
package api

import (
	"net/http"
	"time"

	"user-segmentation-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=fst_user_segments password=18081971 sslmode=disable")
	if err != nil {
		panic("Failed to connect to database")
	}
	// Создание таблицы сегментов, если её еще нет
	db.AutoMigrate(&models.Segment{})
}

// CreateSegment обработчик для создания сегмента
func CreateSegment(c *gin.Context) {
	var segment models.Segment
	if err := c.ShouldBindJSON(&segment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Попытка создать сегмент в базе данных
	if err := db.Create(&segment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create segment"})
		return
	}

	// После успешного создания сегмента, возвращаем ответ
	c.JSON(http.StatusCreated, gin.H{"message": "Segment created successfully"})
}

// DeleteSegment обработчик для удаления сегмента
func DeleteSegment(c *gin.Context) {
	slug := c.Param("slug")

	// Попытка найти сегмент по slug в базе данных
	var segment models.Segment
	if err := db.Where("slug = ?", slug).First(&segment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
		return
	}

	// Попытка удалить сегмент из базы данных
	if err := db.Delete(&segment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete segment"})
		return
	}

	// После успешного удаления сегмента, возвращаем ответ
	c.JSON(http.StatusOK, gin.H{"message": "Segment deleted successfully"})
}

// AddUserToSegment обработчик для добавления пользователя в сегмент
func AddUserToSegment(c *gin.Context) {
	var requestData struct {
		SegmentsToAdd    []string `json:"segments_to_add"`
		SegmentsToRemove []string `json:"segments_to_remove"`
		UserID           int      `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Попытка найти пользователя по ID в базе данных
	var user models.User
	if err := db.First(&user, requestData.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Добавление сегментов
	for _, slugToAdd := range requestData.SegmentsToAdd {
		var segment models.Segment
		if err := db.Where("slug = ?", slugToAdd).First(&segment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
			return
		}

		// Проверка, если пользователь уже состоит в этом сегменте
		var userSegment models.UserSegment
		if err := db.Where("user_id = ? AND segment_id = ?", user.ID, segment.ID).First(&userSegment).Error; err != nil {
			// Создание записи о принадлежности пользователя к сегменту
			userSegment = models.UserSegment{
				UserID:    user.ID,
				SegmentID: segment.ID,
				AddedAt:   time.Now().UTC().Format(time.RFC3339), // текущее время
			}
			db.Create(&userSegment)
		}
	}

	// Удаление сегментов
	for _, slugToRemove := range requestData.SegmentsToRemove {
		var segment models.Segment
		if err := db.Where("slug = ?", slugToRemove).First(&segment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
			return
		}

		// Удаление записи о принадлежности пользователя к сегменту
		db.Where("user_id = ? AND segment_id = ?", user.ID, segment.ID).Delete(&models.UserSegment{})
	}

	// После успешного добавления/удаления сегментов, возвращаем ответ
	c.JSON(http.StatusOK, gin.H{"message": "Segments updated successfully"})
}

// GetActiveSegments обработчик для получения активных сегментов пользователя
func GetActiveSegments(c *gin.Context) {
	userID := c.Query("user_id")

	// Попытка найти пользователя по ID в базе данных
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var userSegments []models.UserSegment
	// Поиск записей о принадлежности пользователя к сегментам
	if err := db.Where("user_id = ?", user.ID).Find(&userSegments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user segments"})
		return
	}

	var activeSegments []string
	// Поиск и добавление активных сегментов
	for _, userSegment := range userSegments {
		var segment models.Segment
		if err := db.First(&segment, userSegment.SegmentID).Error; err != nil {
			continue // Пропустить этот сегмент, если не найден
		}
		activeSegments = append(activeSegments, segment.Slug)
	}

	c.JSON(http.StatusOK, gin.H{"segments": activeSegments})
}
