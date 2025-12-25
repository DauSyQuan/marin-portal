package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. Lấy danh sách Thủy thủ theo Tàu
func GetCrewByShip(c *gin.Context) {
	shipID := c.Param("ship_id")
	var crews []models.Crew

	// Lấy dữ liệu, người mới tạo lên đầu
	if err := database.DB.Where("ship_id = ?", shipID).Order("created_at desc").Find(&crews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, crews)
}

// 2. Thêm mới Thủy thủ
func AddCrew(c *gin.Context) {
	var input models.Crew
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gán giá trị mặc định
	if input.Status == "" { input.Status = "Active" }
	if input.DataPlan == "" { input.DataPlan = "Basic (1GB)" }
	input.DataUsage = 0
	input.CreatedAt = time.Now()

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi Database"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// 3. Cập nhật thông tin (Sửa)
func UpdateCrew(c *gin.Context) {
	id := c.Param("id")
	var input models.Crew

	// Validate dữ liệu gửi lên
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tìm user trong DB
	var crew models.Crew
	if err := database.DB.First(&crew, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User không tồn tại"})
		return
	}

	// Cập nhật các trường
	database.DB.Model(&crew).Updates(models.Crew{
		FullName:    input.FullName,
		Username:    input.Username,
		Rank:        input.Rank,
		DataPlan:    input.DataPlan,
		Status:      input.Status,
		Nationality: input.Nationality,
	})

	c.JSON(http.StatusOK, crew)
}

// 4. Xóa Thủy thủ
func DeleteCrew(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Crew{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi xóa dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đã xóa thành công"})
}