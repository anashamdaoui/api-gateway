package main

import (
	"log"
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routers"

	"github.com/gorilla/mux"
)

func main() {
	// Load the configuration file
	cfg, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the logger with the configured log level
	middleware.InitLogger(cfg.LogLevel)

	// Set up the router with middleware
	router := mux.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.LoggerMiddleware)
	router.Use(middleware.ErrorHandler)

	// Initialize routes
	routers.InitRoutes(router, cfg.RegistryURL)

	// Start the HTTP server
	log.Printf("API Gateway is running on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
