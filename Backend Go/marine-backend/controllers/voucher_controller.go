package controllers

import (
	"fmt"
	"marine-backend/database"
	"marine-backend/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 1. Lấy danh sách Voucher (Có lọc)
func GetVouchers(c *gin.Context) {
	status := c.Query("status")
	assignTo := c.Query("assign_to")

	var vouchers []models.Voucher
	query := database.DB.Order("created_at desc")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if assignTo != "" {
		query = query.Where("assign_to ILIKE ?", "%"+assignTo+"%")
	}

	query.Find(&vouchers)
	c.JSON(http.StatusOK, vouchers)
}

// 2. Tạo Voucher Mới (Random Code)
func CreateVoucher(c *gin.Context) {
	var input models.Voucher
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sinh mã ngẫu nhiên: VOU-12345
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	input.Code = fmt.Sprintf("VOU-%05d", r.Intn(99999))
	
	input.Status = "Unused"
	input.CreatedBy = "admin" // Mặc định hoặc lấy từ Token
	input.CreatedAt = time.Now()
	if input.ValidDays == 0 { input.ValidDays = 30 }

	database.DB.Create(&input)
	c.JSON(http.StatusCreated, input)
}

// 3. Gán Voucher cho Crew
func AssignVoucher(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CrewID   uint   `json:"crew_id"`
		CrewName string `json:"crew_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu sai"})
		return
	}

	// Cập nhật DB
	result := database.DB.Model(&models.Voucher{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    "Assigned",
		"crew_id":   req.CrewID,
		"assign_to": req.CrewName,
	})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voucher không tồn tại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã gán thành công"})
}