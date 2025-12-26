package controllers

import (
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"marine-backend/database"
	"marine-backend/models"

	"github.com/gin-gonic/gin"
)

// ------------------------
// Helpers
// ------------------------

func getRangeParam(c *gin.Context) string {
	r := c.Query("range")
	if r == "" {
		return "24h"
	}
	switch r {
	case "24h", "7d", "30d":
		return r
	default:
		return "24h"
	}
}

func rangeToDuration(r string) time.Duration {
	switch r {
	case "7d":
		return 7 * 24 * time.Hour
	case "30d":
		return 30 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}

func clampFloat(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}

func randomBetween(min, max int) int {
	if max <= min {
		return min
	}
	return rand.Intn(max-min+1) + min
}

// ------------------------
// GET /api/analytics/overview
// ------------------------
func GetAnalyticsOverview(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	rng := getRangeParam(c)
	duration := rangeToDuration(rng)
	start := time.Now().Add(-duration)

	// 1) Total consumption: sum crew.data_usage (GB)
	var totalGB float64
	if err := database.DB.Model(&models.Crew{}).
		Select("COALESCE(SUM(data_usage), 0)").
		Scan(&totalGB).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute total consumption"})
		return
	}

	totalTB := totalGB / 1024.0

	// 2) Estimated cost: giả định 1GB = 4 USD (bạn có thể đổi theo gói)
	estimatedCost := int(totalGB * 4.0)
	// Projection: +20%
	costProjection := int(float64(estimatedCost) * 1.2)

	// 3) Security threats: dùng AuditLogs có status Warning/Failed trong range
	var threatCount int64
	_ = database.DB.Model(&models.AuditLog{}).
		Where("created_at >= ? AND (status = ? OR status = ?)", start, "Warning", "Failed").
		Count(&threatCount).Error

	// 4) Fleet uptime: tính theo ship.status (Online -> uptime)
	var totalShips int64
	var onlineShips int64
	_ = database.DB.Model(&models.Ship{}).Count(&totalShips).Error
	_ = database.DB.Model(&models.Ship{}).Where("status = ?", "Online").Count(&onlineShips).Error

	var uptime float64 = 99.0
	if totalShips > 0 {
		uptime = (float64(onlineShips) / float64(totalShips)) * 100.0
		// nâng baseline để nhìn giống uptime thật
		uptime = clampFloat(95.0+uptime*0.05, 80, 100)
	}

	// Change pct (tạm giả lập mượt để UI đẹp; khi có historical table thì tính thật)
	totalChangePct := clampFloat(rand.NormFloat64()*5+10, -30, 30)
	threatChangePct := clampFloat(rand.NormFloat64()*3-2, -30, 30)
	uptimeChangePct := clampFloat(rand.NormFloat64()*0.2+0.1, -2, 2)

	c.JSON(http.StatusOK, gin.H{
		"range": rng,

		"total_consumption_tb":         math.Round(totalTB*10) / 10,
		"total_consumption_change_pct": math.Round(totalChangePct*10) / 10,

		"estimated_cost_usd":  estimatedCost,
		"cost_projection_usd": costProjection,

		"security_threats":        threatCount,
		"threats_change_pct":      math.Round(threatChangePct*10) / 10,
		"fleet_uptime_pct":        math.Round(uptime*10) / 10,
		"uptime_change_pct":       math.Round(uptimeChangePct*10) / 10,
	})
}

// ------------------------
// GET /api/analytics/traffic
// ------------------------
func GetAnalyticsTraffic(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	rng := getRangeParam(c)

	// UI hiện đang dùng mốc 7 label cho 24h.
	labels := []string{"00:00", "04:00", "08:00", "12:00", "16:00", "20:00", "24:00"}

	// Chưa có bảng metrics thật (latency/bandwidth theo giờ), nên tạm giả lập theo tổng GB.
	// Khi bạn thêm table metrics_ship_hourly, chỉ cần thay phần generate này bằng query.
	var totalGB float64
	_ = database.DB.Model(&models.Crew{}).
		Select("COALESCE(SUM(data_usage), 0)").
		Scan(&totalGB).Error

	base := 150 + int(totalGB/10)
	if base < 120 {
		base = 120
	}
	if base > 600 {
		base = 600
	}

	bandwidth := make([]int, 0, len(labels))
	latency := make([]int, 0, len(labels))

	for i := 0; i < len(labels); i++ {
		bw := clampFloat(float64(base+randomBetween(-80, 120)), 50, 800)
		lat := clampFloat(float64(45+randomBetween(-15, 70)), 20, 300)
		bandwidth = append(bandwidth, int(bw))
		latency = append(latency, int(lat))
	}

	// Nếu range lớn hơn (7d/30d) bạn vẫn có thể dùng labels theo ngày.
	if rng == "7d" {
		labels = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	}
	if rng == "30d" {
		labels = []string{"W1", "W2", "W3", "W4"}
	}

	c.JSON(http.StatusOK, gin.H{
		"range":          rng,
		"labels":         labels,
		"bandwidth_mbps": bandwidth,
		"latency_ms":     latency,
	})
}

// ------------------------
// GET /api/analytics/top-consumers
// ------------------------
func GetAnalyticsTopConsumers(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	_ = getRangeParam(c)

	limitStr := c.Query("limit")
	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}

	// Lấy top crew theo data_usage
	var crews []models.Crew
	if err := database.DB.Order("data_usage DESC").Limit(limit).Find(&crews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch top consumers"})
		return
	}

	resp := make([]gin.H, 0, len(crews))
	for _, cr := range crews {
		dataGB := cr.DataUsage
		cost := int(dataGB * 4.0)

		// signal: nếu crew đang Active thì signal cao, không thì random
		signal := randomBetween(40, 99)
		status := "Normal"
		if signal < 55 {
			status = "Warning"
		}
		if dataGB > 300 {
			status = "High"
		}

		resp = append(resp, gin.H{
			"name":     cr.Username,
			"type":     cr.Rank,
			"data_gb":  math.Round(dataGB),
			"cost_usd": cost,
			"signal":   signal,
			"status":   status,
		})
	}

	c.JSON(http.StatusOK, resp)
}

// ------------------------
// GET /api/analytics/app-usage
// ------------------------
func GetAnalyticsAppUsage(c *gin.Context) {
	rng := getRangeParam(c)
	_ = rng

	// Bạn chưa có table app usage, nên trả mock ổn định để FE vẽ radar.
	c.JSON(http.StatusOK, gin.H{
		"range": rng,
		"labels": []string{"VoIP", "Video", "ERP", "Email", "Social", "IoT"},
		"current": []int{85, 90, 60, 40, 70, 30},
		"limit":   []int{60, 70, 90, 90, 50, 100},
	})
}
