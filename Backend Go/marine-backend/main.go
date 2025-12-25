package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

// --- 1. CẤU HÌNH HỆ THỐNG ---
const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "123"           // <--- Mật khẩu DB của bạn
	DB_NAME     = "marine_portal" // <--- Tên Database
	JWT_KEY     = "bi_mat_khong_the_bat_mi"
)

var db *sql.DB

// --- 2. CÁC STRUCT (QUAN TRỌNG: JSON TAG PHẢI VIẾT THƯỜNG) ---
type Ship struct {
	ID        string  `json:"id"`        // Frontend đọc: ship.id
	Name      string  `json:"name"`      // Frontend đọc: ship.name
	Company   string  `json:"company"`   // Frontend đọc: ship.company
	Type      string  `json:"type"`      // Frontend đọc: ship.type
	IP        string  `json:"ip"`        // Frontend đọc: ship.ip
	Satellite string  `json:"satellite"` // Frontend đọc: ship.satellite
	Beam      string  `json:"beam"`      // Frontend đọc: ship.beam
	Status    string  `json:"status"`    // Frontend đọc: ship.status
	Lat       float64 `json:"lat"`       // Frontend đọc: ship.lat
	Lon       float64 `json:"lon"`       // Frontend đọc: ship.lon
	SNR       float64 `json:"snr"`       // Frontend đọc: ship.snr
}

type User struct {
	Username string
	Password string
	FullName string
	Role     string
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// --- 3. MAIN FUNCTION ---
func main() {
	// A. Kết nối Database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Lỗi kết nối:", err)
	}
	defer db.Close()

	// Kiểm tra kết nối
	if err = db.Ping(); err != nil {
		log.Fatal("❌ Không thể kết nối DB (Sai mật khẩu?):", err)
	}
	fmt.Println("✅ Đã kết nối Database thành công!")

	// B. Khởi tạo dữ liệu & Simulator
	seedAdminUser()
	go startSimulation()

	// C. Định tuyến (Router)
	r := mux.NewRouter()
	
	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/api/ships", getShips).Methods("GET")
	r.HandleFunc("/api/ships", createShip).Methods("POST") // API Thêm tàu
	r.HandleFunc("/api/report/{id}", downloadReport).Methods("GET")

	// D. Cấu hình CORS (Cho phép Frontend gọi vào)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("🚀 Backend Go đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}

// --- 4. CÁC API HANDLERS ---

// API: Đăng nhập
func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)

	var storedUser User
	err := db.QueryRow("SELECT username, password, role FROM users WHERE username=$1", creds.Username).Scan(&storedUser.Username, &storedUser.Password, &storedUser.Role)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Wrong password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": storedUser.Username,
		"role":     storedUser.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(JWT_KEY))

	json.NewEncoder(w).Encode(map[string]string{
		"token":    tokenString,
		"username": storedUser.Username,
		"role":     storedUser.Role,
	})
}

// API: Lấy danh sách tàu
func getShips(w http.ResponseWriter, r *http.Request) {
	// Sắp xếp theo thời gian cập nhật mới nhất
	rows, err := db.Query("SELECT id, name, company, type, ip, satellite, beam, lat, lon, status, snr FROM ships ORDER BY updated_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ships []Ship
	for rows.Next() {
		var s Ship
		if err := rows.Scan(&s.ID, &s.Name, &s.Company, &s.Type, &s.IP, &s.Satellite, &s.Beam, &s.Lat, &s.Lon, &s.Status, &s.SNR); err != nil {
			log.Println("Scan Error:", err)
			continue
		}
		ships = append(ships, s)
	}
	json.NewEncoder(w).Encode(ships)
}

// API: Thêm tàu mới (Add Vessel)
func createShip(w http.ResponseWriter, r *http.Request) {
	var s Ship
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Giá trị mặc định
	if s.Status == "" { s.Status = "Online" }
	if s.SNR == 0 { s.SNR = 12.0 }

	sqlStatement := `
	INSERT INTO ships (id, name, company, type, ip, satellite, beam, status, lat, lon, snr)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	
	_, err := db.Exec(sqlStatement, s.ID, s.Name, s.Company, s.Type, s.IP, s.Satellite, s.Beam, s.Status, s.Lat, s.Lon, s.SNR)
	if err != nil {
		log.Println("Insert Error:", err)
		http.Error(w, "Lỗi Database: Có thể trùng ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// API: Xuất PDF
func downloadReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shipID := vars["id"]

	var s Ship
	err := db.QueryRow("SELECT id, name, company, type, status, snr FROM ships WHERE id=$1", shipID).Scan(&s.ID, &s.Name, &s.Company, &s.Type, &s.Status, &s.SNR)
	if err != nil {
		http.Error(w, "Ship not found", 404)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "MARINE PORTAL - VESSEL REPORT")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Report Date: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Vessel Name: %s", s.Name))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("IMO Number:  %s", s.ID))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Owner:       %s", s.Company))
	pdf.Ln(20)

	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(95, 10, "Current Status", "1", 0, "", true, 0, "")
	pdf.CellFormat(95, 10, s.Status, "1", 1, "", false, 0, "")
	pdf.CellFormat(95, 10, "Signal Quality (SNR)", "1", 0, "", true, 0, "")
	pdf.CellFormat(95, 10, fmt.Sprintf("%.1f dB", s.SNR), "1", 1, "", false, 0, "")

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Report-%s.pdf", s.ID))
	pdf.Output(w)
}

// --- 5. CÁC HÀM PHỤ TRỢ ---

// Simulator: Tự động đổi SNR (đã làm tròn 1 số thập phân)
func startSimulation() {
	for {
		_, err := db.Exec(`
			UPDATE ships 
			SET snr = ROUND((snr + (random() - 0.5))::numeric, 1), 
			    updated_at = NOW()
			WHERE status = 'Online'
		`)
		if err != nil { log.Println("Sim Error:", err) }
		time.Sleep(2 * time.Second)
	}
}

// Seed: Tạo User Admin
func seedAdminUser() {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		bytes, _ := bcrypt.GenerateFromPassword([]byte("123"), 14)
		db.Exec("INSERT INTO users (username, password, full_name, role) VALUES ($1, $2, $3, $4)",
			"admin", string(bytes), "System Administrator", "Admin")
		fmt.Println("👤 Đã tạo user mặc định: admin / 123")
	}
}