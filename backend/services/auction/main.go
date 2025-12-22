package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/handler"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/repository"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auction/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load config
	cfg := config.LoadConfig("auction")
	log := logger.New(logger.Config{
		Level:       "info",
		Development: cfg.Env == "development",
		ServiceName: cfg.ServiceName,
	})

	// Connect to DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database", zap.Error(err))
	}

	// Setup layers
	repo := repository.NewPostgresRepo(db)
	svc := service.NewAuctionService(repo)

	grpcHandler := handler.NewGrpcHandler(svc)
	httpHandler := handler.NewHttpHandler(svc)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatal("failed to listen grpc", zap.Error(err))
		}

		s := grpc.NewServer()
		pb.RegisterAuctionServiceServer(s, grpcHandler)
		reflection.Register(s)

		log.Info("Auction gRPC Service starting on port " + cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve grpc", zap.Error(err))
		}
	}()

	// Start HTTP server
	tm := auth.NewTokenManager(cfg.JWTSecret)
	r := handler.SetupRouter(httpHandler, tm)

	// Graceful shutdown
	go func() {
		log.Info("Auction HTTP Service starting on port " + cfg.HTTPPort)
		// Use standard http.Server
		if err := r.Run(":" + cfg.HTTPPort); err != nil {
			log.Fatal("Failed to run http server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
}
