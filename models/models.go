// models/models.go
package models

import (
    "github.com/jinzhu/gorm"
)

// Segment структура для сегмента
type Segment struct {
    gorm.Model
    Slug string `json:"slug"`
}

// User структура для пользователя
type User struct {
    gorm.Model
    // Добавьте поля, представляющие пользователя (например, имя, электронная почта и т.д.)
}

// UserSegment структура для связи пользователя и сегмента
type UserSegment struct {
    gorm.Model
    UserID    uint   `json:"user_id"`
    SegmentID uint   `json:"segment_id"`
    AddedAt   string `json:"added_at"`
}
