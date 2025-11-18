package main

import (
	"context"
	"log"
	"net"

	rass "rastochnoy/genproto/rastochnoy"
	"rastochnoy/api"
	"rastochnoy/api/handler"
	"rastochnoy/config"
	connectiondb "rastochnoy/connection_db"
	"rastochnoy/service"
	writerastochnoy "rastochnoy/write_rastochnoy"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1️⃣ Config yuklash
	cfg := config.Load()

	// 2️⃣ gRPC server uchun listener yaratish
	listener, err := net.Listen("tcp", cfg.USER_SERVICE)
	if err != nil {
		log.Fatal("❌ gRPC portni tinglashda xato:", err)
	}
	defer listener.Close()

	// 3️⃣ DB ulanish
	db, err := connectiondb.ConnectDB()
	if err != nil {
		log.Fatal("❌ DB ulanishda xato:", err)
	}

	// 4️⃣ Service yaratish
	crud := writerastochnoy.NewRastochnoyRepo(db)
	service := service.NewRastochnoyService(crud)

	// 5️⃣ gRPC serverni ishga tushirish
	server := grpc.NewServer()
	rass.RegisterRastochnoyServer(server, service)
	log.Printf("🚀 gRPC server ishlayapti: %v", listener.Addr())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal("❌ gRPC serverda xato:", err)
		}
	}()

	// 6️⃣ HTTP handler va router
	hand := NewHandler(cfg)
	router := api.RegisterRastochnoyRoutes(hand)

	// 7️⃣ WebSocket (DB37 / READ) uchun reader ulash
	writerastochnoy.InitRastochnoyWS(func(ctx context.Context, req *rass.ReadRastochnoyReq) (*rass.ReadRastochnoyRes, error) {
		return service.ReadRastochnoy(ctx, req)
	})

	// 8️⃣ WebSocket (DB9 / WRITE) uchun reader ulash
	writerastochnoy.InitRastochnoyWSDB9(func(ctx context.Context, req *rass.ReadWriteRastochnoyReq) (*rass.ReadWriteRastochnoyRes, error) {
		return service.ReadWriteRastochnoy(ctx, req)
	})

	// 10 WebSocket (DB33 / WRITE) uchun reader ulash
	writerastochnoy.InitRastochnoyWSdb33(func(ctx context.Context, req *rass.ReadWriteRastochnoyDB33Req) (*rass.ReadWriteRastochnoyDB33Res, error) {
		return service.ReadWriteRastochnoyDB33(ctx, req)
	})

	// 9️⃣ HTTP (REST + WebSocket) serverni ishga tushirish
	log.Printf("🌐 HTTP + WebSocket server ishlayapti: %s", cfg.USER_ROUTER)
	if err := router.Run(cfg.USER_ROUTER); err != nil {
		log.Fatal("❌ HTTP serverda xato:", err)
	}
}

// gRPC client ulash (handler uchun)
func NewHandler(cfg config.Config) *handler.Handler {
	conn, err := grpc.Dial(cfg.USER_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ gRPC client ulanishida xato: %v", err)
	}

	return &handler.Handler{
		Rastochnoy: rass.NewRastochnoyClient(conn),
	}
}
