package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. Lấy cấu hình hiện tại
func GetSettings(c *gin.Context) {
	var config models.SystemConfig
	// Lấy bản ghi đầu tiên, nếu chưa có thì tạo mặc định
	if err := database.DB.First(&config).Error; err != nil {
		config = models.SystemConfig{
			PrimaryLink: "vsat_ka", SdwanMode: "failover", Firewall: true,
			SnrThreshold: 5.0, QuotaWarning: 80, Recipients: "admin@marine.com",
		}
		database.DB.Create(&config)
	}
	c.JSON(http.StatusOK, config)
}

// 2. Cập nhật cấu hình & Ghi Log
func UpdateSettings(c *gin.Context) {
	var input models.SystemConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cập nhật vào DB (Bản ghi ID=1)
	// Sử dụng Model để đảm bảo chỉ update ID 1
	var config models.SystemConfig
	database.DB.First(&config)
	
	// Copy dữ liệu mới đè lên
	input.ID = config.ID 
	if err := database.DB.Save(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lưu cấu hình"})
		return
	}

	// --- GHI AUDIT LOG ---
	// Lấy IP người dùng (Giả lập hoặc lấy từ Header)
	clientIP := c.ClientIP()
	
	logEntry := models.AuditLog{
		User:      "admin", // Thực tế lấy từ JWT Token
		Action:    "Updated System Configuration",
		IPAddress: clientIP,
		Status:    "Success",
		CreatedAt: time.Now(),
	}
	database.DB.Create(&logEntry)

	c.JSON(http.StatusOK, gin.H{"message": "Cấu hình đã lưu & Ghi nhật ký thành công"})
}

// 3. Lấy danh sách Audit Logs
func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	// Lấy 50 log mới nhất
	database.DB.Order("created_at desc").Limit(50).Find(&logs)
	c.JSON(http.StatusOK, logs)
}