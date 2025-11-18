package api

import (
	"rastochnoy/api/handler"
	writerastochnoy "rastochnoy/write_rastochnoy"

	"github.com/gin-gonic/gin"
)

func RegisterRastochnoyRoutes(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	r := router.Group("/rastochnoy")
	{
		r.PUT("/write", h.WriteRastochnoy)
		r.GET("/readwrite", h.ReadWriteRastochnoy)
		r.GET("/read", h.ReadRastochnoy)
		r.GET("/ws", writerastochnoy.WebSocketHandler) // ✅ WebSocket endpoint
		r.GET("/wsdb9",writerastochnoy.WebSocketHandlerDB9)
		r.PUT("/writeRastochnoy_db37", h.WrtieRastochnoy_db37)
		r.PUT("/writedb33", h.WriteRastochnoyDB33)
		r.GET("/readwritedb33", h.ReadWriteRastochnoyDB33)
		r.GET("/wsdb33", writerastochnoy.WebSocketHandlerdb33)
	}

	return router
}
