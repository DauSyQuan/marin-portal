package models

import "time"

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
	Crews     []Crew    `json:"crews" gorm:"foreignKey:ShipID"`
}

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

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"-"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetInput struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	SecretKey   string `json:"secret_key" binding:"required"`
}
type Voucher struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code"`       // Mã Voucher (VOU-XXXX)
	DataPlan  string    `json:"data_plan"`  // Gói cước (1GB, 5GB...)
	Status    string    `json:"status"`     // Unused, Assigned, Used, Expired
	CreatedBy string    `json:"created_by"` // Admin tạo
	
	// Liên kết với Thủy thủ (Crew)
	AssignTo  string    `json:"assign_to"`  // Tên thủy thủ (Lưu tên cho đơn giản)
	CrewID    uint      `json:"crew_id"`    // ID thủy thủ
	
	UsedAt    time.Time `json:"used_at"`
	CreatedAt time.Time `json:"created_at"`
	ValidDays int       `json:"valid_days"` // Hạn sử dụng (ngày)
}