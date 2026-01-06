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

// ... Giữ nguyên hàm GetVouchers ...

// TẠO VOUCHER MỚI
func CreateVoucher(c *gin.Context) {
	var input models.Voucher
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sinh mã
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	input.Code = fmt.Sprintf("VOU-%05d", r.Intn(99999))
	input.Status = "Unused"
	input.CreatedAt = time.Now()

	// Lưu DB
	database.DB.Create(&input)

	// Đẩy xuống MikroTik ngay lập tức (Tạo sẵn user chờ khách nhập)
	shipID := "IMO9562623" // ID tàu đang test
	client, _, err := ConnectToRouter(shipID)
	
	if err == nil {
		defer client.Close()
		// Tạo User: Name=Code, Password=Code, Profile=DataPlan
		_, err := client.Run(
			"/ip/hotspot/user/add",
			"=name=" + input.Code,
			"=password=" + input.Code,
			"=profile=" + input.DataPlan, // Tên gói cước phải khớp với tên Profile
			"=comment=Voucher",
		)
		if err != nil {
			fmt.Println("⚠️ Lỗi tạo Voucher trên Router:", err)
		} else {
			fmt.Println("✅ Đã bắn Voucher vào MikroTik:", input.Code)
		}
	}

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