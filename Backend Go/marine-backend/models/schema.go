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
    ID          uint   `json:"id" gorm:"primaryKey"`
    ShipID      string `json:"ship_id"` // Khóa ngoại trỏ về tàu
    FullName    string `json:"full_name"`
    Rank        string `json:"rank"`        // Chức danh: Captain, Chief Engineer...
    Nationality string `json:"nationality"`
    Status      string `json:"status"`      // Onboard, On Leave
}

type User struct {
	Username string `json:"username" gorm:"primaryKey"`
	Password string `json:"-"` // Không trả password về JSON
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// Struct nhận dữ liệu Login
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetInput struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	SecretKey   string `json:"secret_key" binding:"required"`
}