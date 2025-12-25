package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Lấy danh sách thủy thủ của 1 tàu cụ thể
func GetCrewByShip(c *gin.Context) {
    shipID := c.Param("ship_id")
    var crews []models.Crew
    
    // Tìm thủy thủ có ShipID trùng khớp
    result := database.DB.Where("ship_id = ?", shipID).Find(&crews)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy dữ liệu"})
        return
    }
    c.JSON(http.StatusOK, crews)
}

// Thêm thủy thủ mới
func AddCrew(c *gin.Context) {
    var input models.Crew
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.Status == "" { input.Status = "Onboard" }

    if err := database.DB.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi Database"})
        return
    }
    c.JSON(http.StatusCreated, input)
}

// Xóa thủy thủ (Optional: Sa thải/Rời tàu)
func DeleteCrew(c *gin.Context) {
    id := c.Param("id")
    database.DB.Delete(&models.Crew{}, id)
    c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}