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
	clientsdb33      = make(map[*websocket.Conn]bool)
	clientsMudb33    sync.Mutex
	lastSentDatadb33 []byte
	upgraderdb33     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	repoReaderdb33 func(context.Context, *rass.ReadWriteRastochnoyDB33Req) (*rass.ReadWriteRastochnoyDB33Res, error)
)

// InitRastochnoyWS — DB o‘qish funksiyasini biriktirish
func InitRastochnoyWSdb33(reader func(context.Context, *rass.ReadWriteRastochnoyDB33Req) (*rass.ReadWriteRastochnoyDB33Res, error)) {
	repoReaderdb33 = reader
	go startPollingdb33()
}

// WebSocketHandler — WebSocket endpoint
func WebSocketHandlerdb33(c *gin.Context) {
	conn, err := upgraderdb33.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ WebSocket ulanishida xato:", err)
		return
	}

	clientsMudb33.Lock()
	clientsdb33[conn] = true
	clientsMudb33.Unlock()

	log.Println("✅ Yangi WebSocket client ulandi")

	// Hozirgi ma'lumotni yuborish
	if repoReaderdb33 != nil {
		data, err := repoReaderdb33(context.Background(), &rass.ReadWriteRastochnoyDB33Req{})
		if err == nil {
			msg, _ := json.Marshal(data)
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	// Mijozni kuzatish
	go func() {
		defer func() {
			clientsMudb33.Lock()
			delete(clients, conn)
			clientsMudb33.Unlock()
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
func startPollingdb33() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if repoReaderdb33 == nil {
			continue
		}

		data, err := repoReaderdb33(context.Background(), &rass.ReadWriteRastochnoyDB33Req{})
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
		if bytes.Equal(message, lastSentDatadb33) {
			continue
		}

		lastSentDatadb33 = append([]byte(nil), message...)
		broadcastdb33(message)
	}
}

// broadcast — barcha WebSocket mijozlarga yuborish
func broadcastdb33(msg []byte) {
	clientsMudb33.Lock()
	defer clientsMudb33.Unlock()

	for conn := range clientsdb33 {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			conn.Close()
			delete(clients, conn)
		}
	}
}
