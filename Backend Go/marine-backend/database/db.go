package database

import (
	"fmt"
	"log"
	"marine-backend/models" // Thay marine-backend bằng tên module của bạn trong go.mod
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Chỉ log khi có lỗi
	})

	if err != nil {
		log.Fatal("❌ Kết nối DB thất bại:", err)
	}

	// Tự động tạo bảng nếu chưa có (Migration)
	DB.AutoMigrate(&models.Ship{}, &models.User{}, &models.Crew{}, &models.Voucher{}, &models.BandwidthPlan{})

	// Cấu hình Connection Pool
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("✅ Kết nối Database GORM thành công!")
}