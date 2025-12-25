package controllers

import (
	"fmt"
	"marine-backend/database"
	"marine-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func DownloadReport(c *gin.Context) {
	id := c.Param("id")
	var ship models.Ship

	// Lấy thông tin tàu từ DB
	if err := database.DB.Where("id = ?", id).First(&ship).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy tàu"})
		return
	}

	// Tạo PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "MARINE PORTAL - VESSEL REPORT")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Report Date: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Vessel Name: %s", ship.Name))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("IMO Number:  %s", ship.ID))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Owner:       %s", ship.Company))
	pdf.Ln(20)

	// Vẽ bảng trạng thái
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(95, 10, "Current Status", "1", 0, "", true, 0, "")
	pdf.CellFormat(95, 10, ship.Status, "1", 1, "", false, 0, "")
	
	pdf.CellFormat(95, 10, "Signal Quality (SNR)", "1", 0, "", true, 0, "")
	pdf.CellFormat(95, 10, fmt.Sprintf("%.1f dB", ship.SNR), "1", 1, "", false, 0, "")

	// Thiết lập Header để trình duyệt tải file
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=Report-%s.pdf", ship.ID))
	
	// Xuất file ra luồng ghi của Gin
	err := pdf.Output(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo file PDF"})
	}
}