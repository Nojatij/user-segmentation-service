// handlers/handlers.go
package handlers

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "user-segmentation-service/models"
    "user-segmentation-service/db"
)

// CreateSegment обработчик для создания сегмента
func CreateSegment(c *gin.Context) {
    var segment models.Segment
    if err := c.ShouldBindJSON(&segment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
        return
    }

    err := db.CreateSegment(&segment)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create segment"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Segment created successfully"})
}

// DeleteSegment обработчик для удаления сегмента
func DeleteSegment(c *gin.Context) {
    slug := c.Param("slug")

    err := db.DeleteSegment(slug)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete segment"})
        return
    }

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

    err := db.AddUserToSegments(requestData.UserID, requestData.SegmentsToAdd, requestData.SegmentsToRemove)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user segments"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Segments updated successfully"})
}

// GetActiveSegments обработчик для получения активных сегментов пользователя
func GetActiveSegments(c *gin.Context) {
    userID := c.Query("user_id")

    segments, err := db.GetActiveSegments(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user segments"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"segments": segments})
}
