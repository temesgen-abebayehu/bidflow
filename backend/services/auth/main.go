package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/temesgen-abebayehu/bidflow/backend/common/auth"
	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/handler"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/repository"
	"github.com/temesgen-abebayehu/bidflow/backend/services/auth/internal/service"
)

func main() {
	cfg := config.LoadConfig("auth")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Setup layers
	repo := repository.NewPostgresRepo(db)
	companyRepo := repository.NewCompanyRepo(db)
	tm := auth.NewTokenManager(cfg.JWTSecret)

	authSvc := service.NewAuthService(repo, tm)
	userSvc := service.NewUserService(repo, companyRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	r := SetupRouter(authHandler, userHandler, tm)

	log.Printf("Auth Service starting on port %s", cfg.HTTPPort)
	r.Run(":" + cfg.HTTPPort)
}
