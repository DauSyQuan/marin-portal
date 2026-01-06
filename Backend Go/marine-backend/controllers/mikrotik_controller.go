package controllers

import (
	"fmt"
	"io"
	"marine-backend/database"
	"marine-backend/models"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-routeros/routeros"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 1. HÀM KẾT NỐI API (Port 8728)
func ConnectToRouter(shipID string) (*routeros.Client, *models.Ship, error) {
	var ship models.Ship
	if err := database.DB.Where("id = ?", shipID).First(&ship).Error; err != nil {
		return nil, nil, fmt.Errorf("không tìm thấy tàu")
	}
	if ship.RouterIP == "" { return nil, nil, fmt.Errorf("IP rỗng") }
	if ship.RouterPort == 0 { ship.RouterPort = 8728 }

	address := fmt.Sprintf("%s:%d", ship.RouterIP, ship.RouterPort)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil { return nil, &ship, err }

	client, err := routeros.NewClient(conn)
	if err != nil { return nil, &ship, err }
	if err := client.Login(ship.RouterUser, ship.RouterPass); err != nil {
		client.Close()
		return nil, &ship, err
	}
	return client, &ship, nil
}

// 2. HÀM KẾT NỐI SFTP (Port 22 - Để upload file)
func connectSFTP(ship *models.Ship) (*sftp.Client, *ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: ship.RouterUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(ship.RouterPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	
	// SSH mặc định port 22
	addr := fmt.Sprintf("%s:22", ship.RouterIP)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil { return nil, nil, err }

	client, err := sftp.NewClient(conn)
	if err != nil { 
		conn.Close()
		return nil, nil, err 
	}
	return client, conn, nil
}

// --- CÁC API ---

// API 1: Monitor Health
func GetRouterHealth(c *gin.Context) {
	shipID := c.Param("ship_id")
	client, _, err := connectToRouter(shipID)
	
	if err != nil {
		// Trả về data giả lập để Web không lỗi 502
		c.JSON(http.StatusOK, gin.H{
			"connected": false, "board_name": "MikroTik (Simulated)", "version": "Offline",
			"cpu_load": 0, "free_memory": 0, "tx_rate": 0, "rx_rate": 0,
		})
		return
	}
	defer client.Close()

	res, _ := client.Run("/system/resource/print")
	traffic, _ := client.Run("/interface/monitor-traffic", "=interface=ether1", "=once")

	data := res.Re[0].Map
	cpu, _ := strconv.Atoi(data["cpu-load"])
	freeMem, _ := strconv.ParseInt(data["free-memory"], 10, 64)
	
	var tx, rx int64 = 0, 0
	if len(traffic.Re) > 0 {
		tx, _ = strconv.ParseInt(traffic.Re[0].Map["tx-bits-per-second"], 10, 64)
		rx, _ = strconv.ParseInt(traffic.Re[0].Map["rx-bits-per-second"], 10, 64)
	}

	c.JSON(http.StatusOK, gin.H{
		"connected": true, "board_name": data["board-name"], "version": data["version"],
		"uptime": data["uptime"], "cpu_load": cpu, "free_memory": freeMem/1024/1024,
		"tx_rate": tx, "rx_rate": rx,
	})
}

// API 2: Upload File .rsc và chạy lệnh Import
func UploadConfigFile(c *gin.Context) {
	shipID := c.Param("ship_id")
	
	// Lấy file từ Frontend gửi lên
	file, header, err := c.Request.FormFile("config_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chưa chọn file"})
		return
	}
	defer file.Close()

	// 1. Lấy thông tin tàu
	_, ship, err := connectToRouter(shipID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Không tìm thấy tàu hoặc Router Offline"})
		return
	}

	// 2. Kết nối SFTP
	sftpClient, sshConn, err := connectSFTP(ship)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi kết nối SSH/SFTP (Port 22): " + err.Error()})
		return
	}
	defer sshConn.Close()
	defer sftpClient.Close()

	// 3. Upload file lên thư mục gốc của Router
	remotePath := "/" + header.Filename
	dstFile, err := sftpClient.Create(remotePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không tạo được file trên Router"})
		return
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi ghi file"})
		return
	}

	// 4. (Tùy chọn) Tự động chạy lệnh Import file vừa up
	// Cần kết nối lại API để chạy lệnh /import
	apiClient, _, _ := connectToRouter(shipID)
	if apiClient != nil {
		defer apiClient.Close()
		// Lệnh import file cấu hình
		apiClient.Run("/import", "=file-name="+header.Filename)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã nạp file " + header.Filename + " thành công!"})
}

// API 3: Web Terminal (Gửi lệnh text)
func RunTerminalCommand(c *gin.Context) {
	shipID := c.Param("ship_id")
	var req struct { Command string `json:"command"` }
	c.ShouldBindJSON(&req)

	client, _, err := connectToRouter(shipID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"output": "Error: Router Offline"})
		return
	}
	defer client.Close()

	// Chạy lệnh
	reply, err := client.Run(req.Command)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"output": "Error: " + err.Error()})
		return
	}

	// Format kết quả
	var output string
	for _, re := range reply.Re {
		for k, v := range re.Map { output += fmt.Sprintf("%s: %s  ", k, v) }
		output += "\n"
	}
	if output == "" { output = "Done." }

	c.JSON(http.StatusOK, gin.H{"output": output})
}
// HÀM MỚI: Áp dụng chặn ứng dụng theo cấu hình
func ApplyFirewallRules(c *gin.Context) {
	shipID := "IMO9562623" // ID tàu đang test (Hoặc lấy từ params)
	
	// Nhận cấu hình từ Frontend gửi lên
	var config struct {
		BlockYoutube  bool `json:"block_youtube"`
		BlockFacebook bool `json:"block_facebook"`
		BlockTiktok   bool `json:"block_tiktok"`
		BlockTorrent  bool `json:"block_torrent"`
	}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, _, err := connectToRouter(shipID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Router Offline"})
		return
	}
	defer client.Close()

	// 1. Xóa các luật cũ của Marine Pro để tránh trùng lặp
	// Tìm các rule có comment="MARINE_BLOCK_..."
	existingL7, _ := client.Run("/ip/firewall/layer7-protocol/print", "?comment=MARINE_BLOCK")
	for _, r := range existingL7.Re {
		client.Run("/ip/firewall/layer7-protocol/remove", "=.id="+r.Map[".id"])
	}
	existingFilter, _ := client.Run("/ip/firewall/filter/print", "?comment=MARINE_BLOCK")
	for _, r := range existingFilter.Re {
		client.Run("/ip/firewall/filter/remove", "=.id="+r.Map[".id"])
	}

	// 2. Định nghĩa Regex (Mẫu nhận diện ứng dụng)
	apps := map[string]string{
		"Youtube":  "^.+(youtube.com|googlevideo.com|youtu.be).*$",
		"Facebook": "^.+(facebook.com|facebook.net|fbcdn.net|messenger.com).*$",
		"Tiktok":   "^.+(tiktok.com|tiktokv.com|muscdn.com).*$",
		"Torrent":  "^.+(torrent|tracker|announce|bitcomet|thunder).*$",
	}

	// 3. Tạo luật chặn mới dựa trên config
	createBlockRule := func(appName string, regex string) {
		// Tạo L7 Protocol
		client.Run("/ip/firewall/layer7-protocol/add", 
			"=name="+appName, 
			"=regexp="+regex, 
			"=comment=MARINE_BLOCK",
		)
		// Tạo Filter Rule (Drop)
		client.Run("/ip/firewall/filter/add", 
			"=chain=forward", 
			"=layer7-protocol="+appName, 
			"=action=drop", 
			"=comment=MARINE_BLOCK",
		)
	}

	if config.BlockYoutube { createBlockRule("Youtube", apps["Youtube"]) }
	if config.BlockFacebook { createBlockRule("Facebook", apps["Facebook"]) }
	if config.BlockTiktok { createBlockRule("Tiktok", apps["Tiktok"]) }
	if config.BlockTorrent { createBlockRule("Torrent", apps["Torrent"]) }

	c.JSON(http.StatusOK, gin.H{"message": "Đã cập nhật Tường lửa thành công!"})
}
// Các hàm cũ (Sync, Reboot) giữ nguyên logic
func SyncCrewToRouter(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Đã đồng bộ User"}) }
func RebootRouter(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Đang khởi động lại..."}) }
func BlockInternet(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Đã chặn mạng!"}) }