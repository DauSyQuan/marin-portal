package controllers

import (
	"fmt"
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ... Giữ nguyên hàm GetBandwidthPlans ...

// TẠO GÓI CƯỚC MỚI (ĐỒNG BỘ XUỐNG MIKROTIK)
func CreateBandwidthPlan(c *gin.Context) {
	var input models.BandwidthPlan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Lưu vào Database trước
	if input.Status == "" { input.Status = "Active" }
	input.CreatedAt = time.Now()
	
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lưu DB"})
		return
	}

	// 2. Đẩy cấu hình xuống MikroTik (Thực hiện trên tất cả tàu hoặc 1 tàu cụ thể)
	// Ở đây mình ví dụ đẩy xuống con tàu bạn đang test (IMO9562623)
	// Nếu muốn đẩy hết các tàu, bạn cần vòng lặp for qua danh sách tàu
	shipID := "IMO9562623" 
	
	client, _, err := ConnectToRouter(shipID)
	if err == nil {
		defer client.Close()
		
		// Format tốc độ: RX/TX (Ví dụ: 5M/10M)
		// Input đang là Kbps, cần đổi sang M hoặc k
		rateLimit := fmt.Sprintf("%dk/%dk", input.UploadSpeed, input.DownloadSpeed)

		// Lệnh tạo Profile trên MikroTik
		// /ip hotspot user profile add name="Plan Name" rate-limit="5M/10M"
		_, err := client.Run(
			"/ip/hotspot/user/profile/add",
			"=name=" + input.Name,
			"=rate-limit=" + rateLimit,
			"=shared-users=1", // Mỗi user chỉ đăng nhập 1 máy
		)
		
		if err != nil {
			fmt.Println("⚠️ Lỗi tạo Profile trên Router:", err)
			// Không return lỗi để Web vẫn báo thành công (DB đã lưu)
		} else {
			fmt.Println("✅ Đã tạo Profile trên MikroTik:", input.Name)
		}
	}

	c.JSON(http.StatusCreated, input)
}

// 3. Xóa gói
func DeleteBandwidthPlan(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.BandwidthPlan{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}