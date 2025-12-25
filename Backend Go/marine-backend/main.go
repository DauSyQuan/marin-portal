package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

// --- CONFIG (BƯỚC SAU CÓ THỂ ĐƯA RA ENV) ---
const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "123"
	DB_NAME     = "marine_portal"

	JWT_KEY = "bi_mat_khong_the_bat_mi"
)

var db *sql.DB

// --- STRUCTS ---
type Ship struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Company   string  `json:"company"`
	Type      string  `json:"type"`
	IP        string  `json:"ip"`
	Satellite string  `json:"satellite"`
	Beam      string  `json:"beam"`
	Status    string  `json:"status"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	SNR       float64 `json:"snr"`
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

// --- MAIN ---
func main() {
	// 1) CONNECT DB
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Cannot open DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Cannot ping DB:", err)
	}

	// 2) ENSURE TABLES + SEED ADMIN + SIMULATION
	ensureTables()
	seedAdminUser()
	go startSimulation()

	// 3) ROUTER
	r := mux.NewRouter()
	r.HandleFunc("/api/login", login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/ships", getShips).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/ships", createShip).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/report/{id}", downloadReport).Methods("GET", "OPTIONS")

	// 4) CORS (DEV: allow localhost/127 mọi port)
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:")
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		// Bearer token -> không cần credentials
		AllowCredentials: false,
		MaxAge:           300,
	})

	fmt.Println("🚀 Backend Go đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}

// --- JSON HELPERS ---
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{
		"error":   code,
		"message": message,
	})
}

// --- HANDLERS ---

// API 1: Login
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	creds.Username = strings.TrimSpace(creds.Username)
	creds.Password = strings.TrimSpace(creds.Password)

	if creds.Username == "" || creds.Password == "" {
		writeError(w, http.StatusBadRequest, "MISSING_FIELDS", "Username/password is required")
		return
	}

	var storedUser User
	err := db.QueryRow(
		"SELECT username, password, role FROM users WHERE username=$1",
		creds.Username,
	).Scan(&storedUser.Username, &storedUser.Password, &storedUser.Role)

	// Không leak info user tồn tại hay không
	if err != nil {
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Incorrect username or password")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Incorrect username or password")
		return
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": storedUser.Username,
		"role":     storedUser.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWT_KEY))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "TOKEN_ERROR", "Failed to generate token")
		return
	}

	// ✅ JSON chuẩn cho frontend
	writeJSON(w, http.StatusOK, map[string]any{
		"token": tokenString,
		"user": map[string]string{
			"username": storedUser.Username,
			"role":     storedUser.Role,
		},
	})
}

// API 2: Get ships
func getShips(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	rows, err := db.Query(`
		SELECT id, name, company, type, ip, satellite, beam, lat, lon, status, snr
		FROM ships
		ORDER BY updated_at DESC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	defer rows.Close()

	var ships []Ship
	for rows.Next() {
		var s Ship
		if err := rows.Scan(&s.ID, &s.Name, &s.Company, &s.Type, &s.IP, &s.Satellite, &s.Beam, &s.Lat, &s.Lon, &s.Status, &s.SNR); err != nil {
			continue
		}
		ships = append(ships, s)
	}

	writeJSON(w, http.StatusOK, ships)
}

// API 3: Create ship
func createShip(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var s Ship
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", err.Error())
		return
	}

	s.ID = strings.TrimSpace(s.ID)
	s.Name = strings.TrimSpace(s.Name)
	s.Company = strings.TrimSpace(s.Company)

	if s.ID == "" || s.Name == "" || s.Company == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "id/name/company is required")
		return
	}

	// Default values
	if s.Status == "" {
		s.Status = "Online"
	}
	if s.SNR == 0 {
		s.SNR = 12.0
	}

	sqlStatement := `
		INSERT INTO ships (id, name, company, type, ip, satellite, beam, status, lat, lon, snr)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := db.Exec(sqlStatement, s.ID, s.Name, s.Company, s.Type, s.IP, s.Satellite, s.Beam, s.Status, s.Lat, s.Lon, s.SNR)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			writeError(w, http.StatusConflict, "DUPLICATE_ID", "Ship ID already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "DB_ERROR", "Database error")
		log.Println("DB error:", err)
		return
	}

	writeJSON(w, http.StatusCreated, s)
}

// API 4: Export PDF
func downloadReport(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(r)
	shipID := vars["id"]

	var s Ship
	err := db.QueryRow(
		"SELECT id, name, company, type, status, snr FROM ships WHERE id=$1",
		shipID,
	).Scan(&s.ID, &s.Name, &s.Company, &s.Type, &s.Status, &s.SNR)

	if err != nil {
		http.Error(w, "Ship not found", http.StatusNotFound)
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
	_ = pdf.Output(w)
}

// --- BACKGROUND SIMULATION ---
func startSimulation() {
	for {
		_, err := db.Exec(`
			UPDATE ships
			SET snr = ROUND((snr + (random() - 0.5))::numeric, 1),
			    updated_at = NOW()
			WHERE status = 'Online'
		`)
		if err != nil {
			log.Println("Sim Error:", err)
		}
		time.Sleep(2 * time.Second)
	}
}

// --- DB INIT ---
func ensureTables() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL,
			full_name TEXT DEFAULT '',
			role TEXT NOT NULL DEFAULT 'User'
		);

		CREATE TABLE IF NOT EXISTS ships (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			company TEXT NOT NULL,
			type TEXT DEFAULT '',
			ip TEXT DEFAULT '',
			satellite TEXT DEFAULT '',
			beam TEXT DEFAULT '',
			status TEXT NOT NULL DEFAULT 'Online',
			lat DOUBLE PRECISION DEFAULT 0,
			lon DOUBLE PRECISION DEFAULT 0,
			snr DOUBLE PRECISION DEFAULT 12,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatal("❌ ensureTables error:", err)
	}
}

func seedAdminUser() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Println("seedAdminUser: COUNT users error:", err)
		return
	}

	if count == 0 {
		bytes, _ := bcrypt.GenerateFromPassword([]byte("123"), 14)
		_, err := db.Exec(
			"INSERT INTO users (username, password, full_name, role) VALUES ($1, $2, $3, $4)",
			"admin", string(bytes), "System Administrator", "Admin",
		)
		if err != nil {
			log.Println("seedAdminUser insert error:", err)
			return
		}
		fmt.Println("👤 Đã tạo user mặc định: admin / 123")
	}
}
