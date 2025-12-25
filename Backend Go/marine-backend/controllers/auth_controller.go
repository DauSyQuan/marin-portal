package controllers

import (
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User không tồn tại"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai mật khẩu"})
		return
	}

	// Tạo Token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"username": user.Username,
		"role":     user.Role,
	})
}

func ResetPassword(c *gin.Context) {
	var input models.ResetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.SecretKey != os.Getenv("MASTER_KEY") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sai mã bảo mật hệ thống"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	result := database.DB.Model(&models.User{}).Where("username = ?", input.Username).Update("password", string(hash))
	
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User không tìm thấy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đổi mật khẩu thành công"})
}

// Seed Admin
func SeedAdmin() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("123"), 14)
		admin := models.User{Username: "admin", Password: string(hash), FullName: "System Admin", Role: "Admin"}
		database.DB.Create(&admin)
	}
}