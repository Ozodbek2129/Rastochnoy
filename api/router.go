package api

import (
	"log"
	"rastochnoy/api/handler"
	writerastochnoy "rastochnoy/write_rastochnoy"

	"github.com/gin-gonic/gin"
)

func RegisterRastochnoyRoutes(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	r := router.Group("/rastochnoy")
	{
		r.PUT("/write", h.WriteRastochnoy)
		r.GET("/readwrite", h.ReadWriteRastochnoy)
		r.GET("/read", h.ReadRastochnoy)
		r.GET("/ws", writerastochnoy.WebSocketHandler) // ✅ WebSocket endpoint
		r.GET("/wsdb9", writerastochnoy.WebSocketHandlerDB9)
		r.PUT("/writeRastochnoy_db37", h.WrtieRastochnoy_db37)
		r.PUT("/writedb33", h.WriteRastochnoyDB33)
		r.GET("/readwritedb33", h.ReadWriteRastochnoyDB33)
		r.GET("/wsdb33", writerastochnoy.WebSocketHandlerdb33)
	}

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Cors middleware triggered")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
