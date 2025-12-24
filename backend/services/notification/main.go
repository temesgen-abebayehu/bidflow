package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/event"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/handler"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/repository"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/service"
	"github.com/temesgen-abebayehu/bidflow/backend/services/notification/internal/websocket"
	"go.uber.org/zap"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig("notification")
	log := logger.New(logger.Config{
		Level:       "info",
		Development: cfg.Env == "development",
		ServiceName: cfg.ServiceName,
	})
	defer log.Sync()

	// 2. Connect to Database
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

	// 3. Initialize Components
	repo := repository.NewPostgresRepo(db)
	hub := websocket.NewHub(log)
	svc := service.NewNotificationService(repo, hub, log)
	tokenManager := auth.NewTokenManager(cfg.JWTSecret)

	// 4. Start WebSocket Hub
	go hub.Run()

	// 5. Initialize and Start Kafka Consumer
	kafkaConsumer := kafka.NewConsumer(
		cfg.KafkaBrokers,
		[]string{event.TopicAuctionCreated, event.TopicBidPlaced},
		"notification-service-group",
		log,
	)
	consumer := event.NewNotificationConsumer(kafkaConsumer, svc, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumer.Start(ctx)

	// 6. Setup HTTP Server
	h := handler.NewNotificationHandler(svc, hub, tokenManager, log)
	r := handler.SetupRouter(h, tokenManager)

	// 7. Start Server
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Info("Starting HTTP server", zap.String("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// 8. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	if err := kafkaConsumer.Close(); err != nil {
		log.Error("Failed to close kafka consumer", zap.Error(err))
	}

	log.Info("Server exiting")
}
