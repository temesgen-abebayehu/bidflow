package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
	"github.com/temesgen-abebayehu/bidflow/backend/common/kafka"
	"github.com/temesgen-abebayehu/bidflow/backend/common/logger"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/event"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/handler"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/repository"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig("auth")
	log := logger.New(logger.Config{
		Level:       "info",
		Development: cfg.Env == "development",
		ServiceName: cfg.ServiceName,
	})

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	// Setup layers
	repo := repository.NewPostgresRepo(db)
	companyRepo := repository.NewCompanyRepo(db)
	tm := auth.NewTokenManager(cfg.JWTSecret)

	// Kafka Producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, log)
	defer kafkaProducer.Close()
	eventProducer := event.NewKafkaEventProducer(kafkaProducer)

	authSvc := service.NewAuthService(repo, tm, eventProducer)
	userSvc := service.NewUserService(repo, companyRepo, eventProducer)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	r := SetupRouter(authHandler, userHandler, tm)

	log.Info("Auth Service starting on port " + cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("Failed to run server", zap.Error(err))
	}
}
