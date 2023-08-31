package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"user-segmentation-service/api"
)

func TestCreateSegment(t *testing.T) {
	r := gin.Default()
	r.POST("/segments", api.CreateSegment)

	payload := `{"slug":"AVITO_TEST_SEGMENT"}`
	req, _ := http.NewRequest("POST", "/segments", strings.NewReader(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDeleteSegment(t *testing.T) {
	r := gin.Default()
	r.DELETE("/segments/:slug", api.DeleteSegment)

	// Создаем DELETE-запрос
	req, _ := http.NewRequest("DELETE", "/segments/AVITO_TEST_SEGMENT", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddUserToSegment(t *testing.T) {
	r := gin.Default()
	r.POST("/segments/:user_id", api.AddUserToSegment)

	// Создаем POST-запрос с JSON-данными
	payload := `{"add_segments":["AVITO_TEST_SEGMENT"], "remove_segments":["AVITO_ANOTHER_SEGMENT"]}`
	req, _ := http.NewRequest("POST", "/segments/1000", strings.NewReader(payload))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetActiveSegments(t *testing.T) {
	r := gin.Default()
	r.GET("/segments/user/:user_id", api.GetActiveSegments)

	// Создаем GET-запрос
	req, _ := http.NewRequest("GET", "/segments/user/1000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}


// Добавьте другие тесты для других функций
