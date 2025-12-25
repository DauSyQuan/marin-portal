package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Lấy danh sách tàu
func GetShips(c *gin.Context) {
	var ships []models.Ship
	database.DB.Order("updated_at desc").Find(&ships)
	c.JSON(http.StatusOK, ships)
}

// Thêm tàu mới
func CreateShip(c *gin.Context) {
	var input models.Ship
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	if input.Status == "" { input.Status = "Online" }
	if input.SNR == 0 { input.SNR = 12.0 }
	input.UpdatedAt = time.Now()

	result := database.DB.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID đã tồn tại hoặc lỗi DB"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// Chạy ngầm Simulator (Dùng GORM)
func StartSimulation() {
	for {
		// Update SNR ngẫu nhiên cho tàu Online
		database.DB.Exec("UPDATE ships SET snr = ROUND((snr + (random() - 0.5))::numeric, 1), updated_at = NOW() WHERE status = 'Online'")
		time.Sleep(2 * time.Second)
	}
}