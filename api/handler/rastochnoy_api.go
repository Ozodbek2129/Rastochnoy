package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	rass "rastochnoy/genproto/rastochnoy"

	"github.com/gin-gonic/gin"
)

// Rastochnoy uchun handler funksiyalar

// ✅ WriteRastochnoy
// Bazaga qiymat yozish uchun API (update rastochnoy_write)
func (h *Handler) WriteRastochnoy(c *gin.Context) {
	var req rass.WriteRastochnoyReq

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Println("Bodydan ma'lumot olishda xatolik:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.Rastochnoy.WriteRastochnoy(context.Background(), &req)
	if err != nil {
		log.Println("WriteRastochnoy RPC da xatolik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✅ ReadWriteRastochnoy
// rastochnoy_write jadvalidan o‘qish uchun API
func (h *Handler) ReadWriteRastochnoy(c *gin.Context) {
	var req rass.ReadWriteRastochnoyReq

	resp, err := h.Rastochnoy.ReadWriteRastochnoy(context.Background(), &req)
	if err != nil {
		log.Println("ReadWriteRastochnoy RPC da xatolik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✅ ReadRastochnoy
// rastochnoy_read jadvalidan o‘qish uchun API
func (h *Handler) ReadRastochnoy(c *gin.Context) {
	var req rass.ReadRastochnoyReq

	resp, err := h.Rastochnoy.ReadRastochnoy(context.Background(), &req)
	if err != nil {
		log.Println("BeadRastochnoy RPC da xatolik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✅ WrtieRastochnoy_db37
// Bazaga qiymat yozish uchun API (update rastochnoy_write)
func (h *Handler) WrtieRastochnoy_db37(c *gin.Context) {
	var req rass.WrtieRastochnoyDb37Req

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Println("Bodydan ma'lumot olishda xatolik:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.Rastochnoy.WrtieRastochnoyDb37(context.Background(), &req)
	if err != nil {
		log.Println("WrtieRastochnoy_db37 RPC da xatolik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✅ WriteRastochnoyDB33
// Bazaga qiymat yozish uchun API (update rastochnoy_writedb33)
func (h *Handler) WriteRastochnoyDB33(c *gin.Context) {
	var req rass.WriteRastochnoydb33Req

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Println("Bodydan ma'lumot olishda xatolik:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.Rastochnoy.WriteRastochnoyDB33(context.Background(), &req)
	if err != nil {
		log.Println("WriteRastochnoy RPC da xatolik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ✅ ReadWriteRastochnoyDB33
// rastochnoy_writedb33 jadvalidan o‘qish uchun API
func (h *Handler) ReadWriteRastochnoyDB33(c *gin.Context) {
	var req rass.ReadWriteRastochnoyDB33Req

	resp, err := h.Rastochnoy.ReadWriteRastochnoyDB33(context.Background(), &req)
	if err != nil {
		log.Println("ReadWriteRastochnoy RPC da xatolik:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}