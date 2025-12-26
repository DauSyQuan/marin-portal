package routes

import (
	"marine-backend/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Cấu hình CORS (Cho phép Frontend truy cập)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Group API
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/reset-password", controllers.ResetPassword)
		api.GET("/ships", controllers.GetShips)
		api.POST("/ships", controllers.CreateShip)
		api.GET("/ships/:ship_id/crew", controllers.GetCrewByShip) // Lấy crew theo tàu
    	api.POST("/crew", controllers.AddCrew)                     // Thêm crew
    	api.DELETE("/crew/:id", controllers.DeleteCrew)            // Xóa crew
		// PDF Report (Bạn có thể copy logic PDF vào controller riêng sau)
		api.GET("/report/:id", controllers.DownloadReport) 
		api.PUT("/crew/:id", controllers.UpdateCrew)
		api.GET("/usage-report", controllers.GetMonthlyUsage)
		api.GET("/online-users", controllers.GetOnlineUsers)
    	api.POST("/online-users/:username/kick", controllers.KickUser)
		api.GET("/vouchers", controllers.GetVouchers)
		api.POST("/vouchers", controllers.CreateVoucher)
		api.PUT("/vouchers/:id/assign", controllers.AssignVoucher) // API Gán
		api.GET("/bandwidth-plans", controllers.GetBandwidthPlans)
		api.POST("/bandwidth-plans", controllers.CreateBandwidthPlan)
		api.DELETE("/bandwidth-plans/:id", controllers.DeleteBandwidthPlan)
		api.GET("/ships/:ship_id/router/stats", controllers.GetRouterHealth)
		api.POST("/ships/:ship_id/router/sync", controllers.SyncCrewToRouter)
		api.POST("/ships/:ship_id/router/reboot", controllers.RebootRouter)
		api.GET("/settings", controllers.GetSettings)
		api.PUT("/settings", controllers.UpdateSettings)
		api.GET("/audit-logs", controllers.GetAuditLogs)
		// Analytics
		api.GET("/analytics/overview", controllers.GetAnalyticsOverview)
		api.GET("/analytics/traffic", controllers.GetAnalyticsTraffic)
		api.GET("/analytics/top-consumers", controllers.GetAnalyticsTopConsumers)
		api.GET("/analytics/app-usage", controllers.GetAnalyticsAppUsage)

	}

	return r
}