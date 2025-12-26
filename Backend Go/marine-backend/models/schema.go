package models

import "time"

// 1. Tàu (Ship) - Đã có thông tin Router
type Ship struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Company   string    `json:"company"`
	Type      string    `json:"type"`
	IP        string    `json:"ip"`
	Satellite string    `json:"satellite"`
	Beam      string    `json:"beam"`
	Status    string    `json:"status"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	SNR       float64   `json:"snr"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Thông tin quản lý Router MikroTik
	RouterIP   string `json:"router_ip"`
	RouterPort int    `json:"router_port"`
	RouterUser string `json:"router_user"`
	RouterPass string `json:"-"` // Không trả về JSON

	Crews     []Crew    `json:"crews" gorm:"foreignKey:ShipID"`
}

// 2. Thủy thủ / ICT User (Crew)
type Crew struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ShipID      string    `json:"ship_id"`
	FullName    string    `json:"full_name"`
	Rank        string    `json:"rank"`
	Nationality string    `json:"nationality"`
	Username    string    `json:"username"`
	DataPlan    string    `json:"data_plan"`
	DataUsage   float64   `json:"data_usage"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// 3. Voucher
type Voucher struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code"`
	DataPlan  string    `json:"data_plan"`
	Status    string    `json:"status"`
	CreatedBy string    `json:"created_by"`
	AssignTo  string    `json:"assign_to"`
	CrewID    uint      `json:"crew_id"`
	UsedAt    time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`
	ValidDays int       `json:"valid_days"`
}

// 4. Gói Băng thông (BandwidthPlan)
type BandwidthPlan struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name"`
	UploadSpeed    int       `json:"upload_speed"`
	DownloadSpeed  int       `json:"download_speed"`
	BurstLimit     string    `json:"burst_limit"`
	BurstThreshold string    `json:"burst_threshold"`
	BurstTime      string    `json:"burst_time"`
	Priority       int       `json:"priority"`
	LimitAt        string    `json:"limit_at"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// 5. Admin User
type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// 6. Input Structs
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetInput struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	SecretKey   string `json:"secret_key" binding:"required"`
}