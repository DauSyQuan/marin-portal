package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. Lấy danh sách
func GetBandwidthPlans(c *gin.Context) {
	var plans []models.BandwidthPlan
	database.DB.Order("created_at desc").Find(&plans)
	c.JSON(http.StatusOK, plans)
}

// 2. Tạo gói mới
func CreateBandwidthPlan(c *gin.Context) {
	var input models.BandwidthPlan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status == "" { input.Status = "Active" }
	input.CreatedAt = time.Now()

	database.DB.Create(&input)
	c.JSON(http.StatusCreated, input)
}

// 3. Xóa gói
func DeleteBandwidthPlan(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.BandwidthPlan{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}