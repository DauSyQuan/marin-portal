package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Struct trả về cho Frontend (Report Row)
type UsageReport struct {
	ID        uint    `json:"id"`
	Username  string  `json:"username"`
	Download  float64 `json:"download"`   // MB
	Upload    float64 `json:"upload"`     // MB
	TotalData float64 `json:"total_data"` // MB
	TimeUsed  int64   `json:"time_used"`  // Seconds
	Status    string  `json:"status"`
}

// Struct nhận Filter từ Frontend
type UsageFilter struct {
	Username  string `form:"username"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

func GetMonthlyUsage(c *gin.Context) {
	var filter UsageFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Lấy danh sách Crew (Có thể lọc theo username)
	var crews []models.Crew
	query := database.DB.Model(&models.Crew{})
	if filter.Username != "" {
		query = query.Where("username ILIKE ?", "%"+filter.Username+"%")
	}
	query.Find(&crews)

	// 2. Tính toán dữ liệu báo cáo (Simulated Aggregation)
	var reports []UsageReport

	// Tính số ngày trong khoảng thời gian chọn (để fake dữ liệu cho hợp lý)
	days := 1.0
	if t1, err := time.Parse("2006-01-02", filter.StartDate); err == nil {
		if t2, err := time.Parse("2006-01-02", filter.EndDate); err == nil {
			days = t2.Sub(t1).Hours() / 24
			if days < 1 { days = 1 }
		}
	}

	for _, crew := range crews {
		// Giả lập: Mỗi ngày dùng khoảng 50-200MB
		// Dùng ID để seed random giúp số liệu ổn định khi refresh
		r := rand.New(rand.NewSource(int64(crew.ID) * int64(days))) 
		
		dl := (50.0 + r.Float64()*150.0) * days
		ul := dl * 0.2 // Upload bằng 20% Download
		
		// Thời gian online (Giả sử 2-5 tiếng/ngày)
		timeUsed := int64((2.0 + r.Float64()*3.0) * 3600 * days)

		reports = append(reports, UsageReport{
			ID:        crew.ID,
			Username:  crew.Username,
			Download:  dl,
			Upload:    ul,
			TotalData: dl + ul,
			TimeUsed:  timeUsed,
			Status:    crew.Status,
		})
	}

	c.JSON(http.StatusOK, reports)
}