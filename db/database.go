// db/database.go
package db

import (
	"time"
	"user-segmentation-service/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

// InitDB инициализация базы данных
func InitDB() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=fst_user_segments password=18081971 sslmode=disable")
	if err != nil {
		panic("Failed to connect to database")
	}
	// Создание таблицы сегментов, если её еще нет
	db.AutoMigrate(&models.Segment{}, &models.UserSegment{})
}

// CreateSegment создание сегмента в базе данных
func CreateSegment(segment *models.Segment) error {
	return db.Create(segment).Error
}

// DeleteSegment удаление сегмента из базы данных по slug
func DeleteSegment(slug string) error {
	return db.Where("slug = ?", slug).Delete(&models.Segment{}).Error
}

// AddUserToSegments добавление/удаление пользователя из сегментов
func AddUserToSegments(userID int, segmentsToAdd, segmentsToRemove []string) error {
	// Найти пользователя в базе данных
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	// Добавить сегменты
	for _, slugToAdd := range segmentsToAdd {
		var segment models.Segment
		if err := db.Where("slug = ?", slugToAdd).First(&segment).Error; err != nil {
			continue
		}

		// Проверить, существует ли уже запись о принадлежности
		var userSegment models.UserSegment
		if err := db.Where("user_id = ? AND segment_id = ?", user.ID, segment.ID).First(&userSegment).Error; err != nil {
			// Создать новую запись
			userSegment = models.UserSegment{
				UserID:    user.ID,
				SegmentID: segment.ID,
				AddedAt:   time.Now().UTC().Format(time.RFC3339),
			}
			db.Create(&userSegment)
		}
	}

	// Удалить сегменты
	for _, slugToRemove := range segmentsToRemove {
		var segment models.Segment
		if err := db.Where("slug = ?", slugToRemove).First(&segment).Error; err != nil {
			continue
		}

		// Удалить запись о принадлежности
		db.Where("user_id = ? AND segment_id = ?", user.ID, segment.ID).Delete(&models.UserSegment{})
	}

	return nil
}

// GetActiveSegments получение активных сегментов пользователя
func GetActiveSegments(userID string) ([]string, error) {
	var segments []string

	// Найти записи о принадлежности пользователя к сегментам
	var userSegments []models.UserSegment
	if err := db.Where("user_id = ?", userID).Find(&userSegments).Error; err != nil {
		return segments, err
	}

	// Для каждой записи о принадлежности, найти информацию о сегменте
	for _, userSegment := range userSegments {
		var segment models.Segment
		if err := db.First(&segment, userSegment.SegmentID).Error; err != nil {
			continue // Пропустить этот сегмент, если не найден
		}
		segments = append(segments, segment.Slug)
	}

	return segments, nil
}
