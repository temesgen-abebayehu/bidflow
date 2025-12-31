package main

import (
	"log"

	"github.com/temesgen-abebayehu/bidflow/backend/common/config"
	"github.com/temesgen-abebayehu/bidflow/backend/services/api-gateway/routes"
)

func main() {
	cfg := config.LoadConfig("api-gateway")

	r := routes.SetupRouter(cfg)

	log.Printf("API Gateway starting on port %s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
