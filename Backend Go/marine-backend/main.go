package main

import (
	"log"
	"marine-backend/controllers"
	"marine-backend/database"
	"marine-backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Env
	godotenv.Load()

	// 2. Connect DB
	database.Connect()

	// 3. Init Data & Workers
	controllers.SeedAdmin()
	go controllers.StartSimulation()

	// 4. Start Server (Gin)
	r := routes.SetupRouter()
	
	log.Println("🚀 Backend Gin running on port 8080")
	r.Run(":8080")
}