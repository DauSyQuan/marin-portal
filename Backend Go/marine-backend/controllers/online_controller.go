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

// Struct trả về dữ liệu hiển thị (Không lưu DB)
type OnlineSession struct {
	Username  string  `json:"username"`
	NasIP     string  `json:"nas_ip"`     // IP thiết bị mạng
	FramedIP  string  `json:"framed_ip"`  // IP cấp cho khách
	Uptime    string  `json:"uptime"`     // Thời gian online (String format)
	Upload    float64 `json:"upload"`     // MB
	Download  float64 `json:"download"`   // MB
	Total     float64 `json:"total"`      // MB
}

// API: Lấy danh sách Online
func GetOnlineUsers(c *gin.Context) {
	// 1. Lấy danh sách Crew đang Active
	var activeCrews []models.Crew
	database.DB.Where("status = ?", "Active").Find(&activeCrews)

	var sessions []OnlineSession
	
	// 2. Giả lập thông số mạng cho từng người
	for i, crew := range activeCrews {
		// Dùng ID để cố định chỉ số ngẫu nhiên (để F5 không bị nhảy số quá nhiều)
		r := rand.New(rand.NewSource(time.Now().Unix() + int64(i))) 

		// Giả lập IP
		framedIP := fmt.Sprintf("192.168.6.%d/32", 100+i)
		
		// Giả lập Traffic (MB)
		down := 50.0 + r.Float64()*500.0
		up := down * 0.15
		
		// Giả lập Uptime (Ví dụ: 0d 2h 30m 10s)
		hours := r.Intn(24)
		mins := r.Intn(60)
		secs := r.Intn(60)
		uptime := fmt.Sprintf("0d %dh %dm %ds", hours, mins, secs)

		sessions = append(sessions, OnlineSession{
			Username: crew.Username,
			NasIP:    "172.29.20.233",
			FramedIP: framedIP,
			Uptime:   uptime,
			Upload:   up,
			Download: down,
			Total:    up + down,
		})
	}

	c.JSON(http.StatusOK, sessions)
}

// API: Kick User (Ngắt kết nối)
func KickUser(c *gin.Context) {
	username := c.Param("username")
	
	// Trong thực tế: Gửi lệnh CoA (Change of Authorization) tới NAS
	// Ở đây: Giả lập thành công
	fmt.Printf("Admin kicked user: %s\n", username)

	c.JSON(http.StatusOK, gin.H{"message": "Đã ngắt kết nối user " + username})
}