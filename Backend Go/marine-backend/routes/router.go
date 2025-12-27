package routes

import (
	"marine-backend/controllers"
	"marine-backend/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Group
	api := r.Group("/api")

	// ✅ Public endpoints (no auth)
	api.POST("/login", controllers.Login)
	api.POST("/reset-password", controllers.ResetPassword)

	// ✅ Protected endpoints (require JWT)
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/ships", controllers.GetShips)
		protected.POST("/ships", controllers.CreateShip)

		protected.GET("/ships/:ship_id/crew", controllers.GetCrewByShip)
		protected.POST("/crew", controllers.AddCrew)
		protected.DELETE("/crew/:id", controllers.DeleteCrew)
		protected.PUT("/crew/:id", controllers.UpdateCrew)

		protected.GET("/report/:id", controllers.DownloadReport)

		protected.GET("/usage-report", controllers.GetMonthlyUsage)

		protected.GET("/online-users", controllers.GetOnlineUsers)
		protected.POST("/online-users/:username/kick", controllers.KickUser)

		protected.GET("/vouchers", controllers.GetVouchers)
		protected.POST("/vouchers", controllers.CreateVoucher)
		protected.PUT("/vouchers/:id/assign", controllers.AssignVoucher)

		protected.GET("/bandwidth-plans", controllers.GetBandwidthPlans)
		protected.POST("/bandwidth-plans", controllers.CreateBandwidthPlan)
		protected.DELETE("/bandwidth-plans/:id", controllers.DeleteBandwidthPlan)

		protected.GET("/ships/:ship_id/router/stats", controllers.GetRouterHealth)
		protected.POST("/ships/:ship_id/router/sync", controllers.SyncCrewToRouter)
		protected.POST("/ships/:ship_id/router/reboot", controllers.RebootRouter)

		protected.GET("/settings", controllers.GetSettings)
		protected.PUT("/settings", controllers.UpdateSettings)
		protected.GET("/audit-logs", controllers.GetAuditLogs)

		// Analytics
		protected.GET("/analytics/overview", controllers.GetAnalyticsOverview)
		protected.GET("/analytics/traffic", controllers.GetAnalyticsTraffic)
		protected.GET("/analytics/top-consumers", controllers.GetAnalyticsTopConsumers)
		protected.GET("/analytics/app-usage", controllers.GetAnalyticsAppUsage)
	}

	return r
}
