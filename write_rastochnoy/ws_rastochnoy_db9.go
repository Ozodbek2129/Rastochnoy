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
	clientsDB9      = make(map[*websocket.Conn]bool)
	clientsMuDB9    sync.Mutex
	lastSentDataDB9 []byte
	upgraderDB9     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	repoReaderDB9 func(context.Context, *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error)
)

// InitRastochnoyWSDB9 — DB9 uchun o‘qish funksiyasini biriktirish
func InitRastochnoyWSDB9(reader func(context.Context, *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error)) {
	repoReaderDB9 = reader
	go startPollingDB9()
}

// WebSocketHandlerDB9 — DB9 uchun WebSocket endpoint
func WebSocketHandlerDB9(c *gin.Context) {
	conn, err := upgraderDB9.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket ulanishida xato:", err)
		return
	}

	clientsMuDB9.Lock()
	clientsDB9[conn] = true
	clientsMuDB9.Unlock()

	log.Println("✅ Yangi WebSocket (DB9) client ulandi")

	// Dastlabki ma'lumotni yuborish
	if repoReaderDB9 != nil {
		data, err := repoReaderDB9(context.Background(), &rass.ReadWriteRastochnoyReq{})
		if err == nil {
			msg, _ := json.Marshal(data)
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	// Client uzilganda tozalash
	go func() {
		defer func() {
			clientsMuDB9.Lock()
			delete(clientsDB9, conn)
			clientsMuDB9.Unlock()
			conn.Close()
			log.Println("⚠️ DB9 client uzildi")
		}()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}

// startPollingDB9 — DB9 ni 0.1 soniyada tekshiradi, o‘zgarsa yuboradi
func startPollingDB9() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if repoReaderDB9 == nil {
			continue
		}

		data, err := repoReaderDB9(context.Background(), &rass.ReadWriteRastochnoyReq{})
		if err != nil {
			log.Println("Polling (DB9) xato:", err)
			continue
		}

		message, err := json.Marshal(data)
		if err != nil {
			log.Println("JSON marshal (DB9) xatosi:", err)
			continue
		}

		// 🔍 faqat o‘zgarish bo‘lsa yubor
		if bytes.Equal(message, lastSentDataDB9) {
			continue
		}

		lastSentDataDB9 = append([]byte(nil), message...)
		broadcastDB9(message)
	}
}

// broadcastDB9 — barcha WebSocket (DB9) mijozlarga yuborish
func broadcastDB9(msg []byte) {
	clientsMuDB9.Lock()
	defer clientsMuDB9.Unlock()

	for conn := range clientsDB9 {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("❌ DB9 clientga yuborishda xato:", err)
			conn.Close()
			delete(clientsDB9, conn)
		}
	}
}
