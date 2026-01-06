package controllers

import (
	"marine-backend/models" // Cần import models để dùng struct Ship
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 1. LẤY DANH SÁCH ONLINE THẬT
func GetOnlineUsers(c *gin.Context) {
	// Lấy ID tàu từ Query hoặc mặc định test
	shipID := "IMO9562623"

	client, _, err := ConnectToRouter(shipID)
	if err != nil {
		// Mất kết nối -> Trả về rỗng
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	defer client.Close()

	// Lệnh lấy danh sách đang online
	// /ip hotspot active print
	res, err := client.Run("/ip/hotspot/active/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var sessions []gin.H
	for _, re := range res.Re {
		data := re.Map
		
		// Chuyển đổi Bytes sang MB
		bytesIn, _ := strconv.ParseFloat(data["bytes-in"], 64)
		bytesOut, _ := strconv.ParseFloat(data["bytes-out"], 64)
		totalMB := (bytesIn + bytesOut) / 1024 / 1024

		sessions = append(sessions, gin.H{
			"username":  data["user"],
			"nas_ip":    data["server"], // Tên Hotspot Server
			"framed_ip": data["address"],
			"uptime":    data["uptime"],
			"upload":    bytesIn / 1024 / 1024,
			"download":  bytesOut / 1024 / 1024,
			"total":     totalMB,
			"id":        data[".id"], // ID dùng để Kick
		})
	}

	c.JSON(http.StatusOK, sessions)
}

// 2. KICK USER THẬT
func KickUser(c *gin.Context) {
	username := c.Param("username") // Thực tế MikroTik cần .id để remove active
	shipID := "IMO9562623"

	client, _, err := ConnectToRouter(shipID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Router Offline"})
		return
	}
	defer client.Close()

	// Tìm ID của session đang active dựa vào username
	findRes, _ := client.Run("/ip/hotspot/active/print", "?user="+username)
	if len(findRes.Re) > 0 {
		id := findRes.Re[0].Map[".id"]
		// Lệnh Kick: /ip hotspot active remove .id=...
		client.Run("/ip/hotspot/active/remove", "=.id="+id)
		c.JSON(http.StatusOK, gin.H{"message": "Đã Kick thành công!"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User không còn online"})
	}
}