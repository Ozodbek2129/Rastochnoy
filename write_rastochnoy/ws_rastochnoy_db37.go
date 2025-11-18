package writerastochnoy

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	rass "rastochnoy/genproto/rastochnoy"
)

var (
	clients      = make(map[*websocket.Conn]bool)
	clientsMu    sync.Mutex
	lastSentData []byte
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	repoReader func(context.Context, *rass.ReadRastochnoyReq) (*rass.ReadRastochnoyRes, error)
)

// InitRastochnoyWS — DB o‘qish funksiyasini biriktirish
func InitRastochnoyWS(reader func(context.Context, *rass.ReadRastochnoyReq) (*rass.ReadRastochnoyRes, error)) {
	repoReader = reader
	go startPolling()
}

// WebSocketHandler — WebSocket endpoint
func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket ulanishida xato:", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()

	log.Println("✅ Yangi WebSocket client ulandi")

	// Hozirgi ma'lumotni yuborish
	if repoReader != nil {
		data, err := repoReader(context.Background(), &rass.ReadRastochnoyReq{})
		if err == nil {
			msg, _ := json.Marshal(data)
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	// Mijozni kuzatish
	go func() {
		defer func() {
			clientsMu.Lock()
			delete(clients, conn)
			clientsMu.Unlock()
			conn.Close()
			log.Println("⚠️ Client uzildi")
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

// startPolling — DBni 0.1 soniyada tekshiradi
func startPolling() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if repoReader == nil {
			continue
		}

		data, err := repoReader(context.Background(), &rass.ReadRastochnoyReq{})
		if err != nil {
			log.Println("Pollingda xato:", err)
			continue
		}

		message, err := json.Marshal(data)
		if err != nil {
			log.Println("JSON marshal xatosi:", err)
			continue
		}

		// 🔍 Faqat o‘zgarish bo‘lsa yubor
		if bytes.Equal(message, lastSentData) {
			continue
		}

		lastSentData = append([]byte(nil), message...)
		broadcast(message)
	}
}

// broadcast — barcha WebSocket mijozlarga yuborish
func broadcast(msg []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("❌ Clientga yuborishda xato:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
