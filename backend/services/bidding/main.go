package main

import (
	"database/sql"
	"fmt"
	"net"
	"os"

	_ "github.com/lib/pq"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	pb "github.com/temesgen-abebayehu/bidflow/backend/proto/pb"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/event"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/handler"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/repository"
	"github.com/temesgen-abebayehu/bidflow/backend/services/bidding/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.LoadConfig("bidding")
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

	// Connect to Auction Service
	auctionSvcURL := os.Getenv("AUCTION_SERVICE_URL")
	if auctionSvcURL == "" {
		auctionSvcURL = "localhost:50051" // Default
	}
	conn, err := grpc.NewClient(auctionSvcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to auction service", zap.Error(err))
	}
	defer conn.Close()
	auctionClient := service.NewAuctionClient(conn)

	// Setup Kafka
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, log)
	defer kafkaProducer.Close()
	eventProducer := event.NewKafkaEventProducer(kafkaProducer)

	// Setup Layers
	repo := repository.NewPostgresRepo(db)
	svc := service.NewBiddingService(repo, eventProducer, auctionClient)

	// Handlers
	httpHandler := handler.NewHttpHandler(svc)
	grpcHandler := handler.NewGrpcHandler(svc)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatal("failed to listen grpc", zap.Error(err))
		}

		s := grpc.NewServer()
		pb.RegisterBiddingServiceServer(s, grpcHandler)
		reflection.Register(s)

		log.Info("Bidding gRPC Service starting on port " + cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve grpc", zap.Error(err))
		}
	}()

	// Start HTTP server
	tm := auth.NewTokenManager(cfg.JWTSecret)
	r := handler.SetupRouter(httpHandler, tm)

	log.Info("Bidding HTTP Service starting on port " + cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("Failed to run server", zap.Error(err))
	}
}
