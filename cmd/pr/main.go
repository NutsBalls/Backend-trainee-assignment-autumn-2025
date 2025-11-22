package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpDelivery "github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/delivery/http"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/repository/postgres"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/internal/pr/usecase/service"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/pkg/config"
	"github.com/NutsBalls/Backend-trainee-assignment-autumn-2025/pkg/database"
)

func main() {
	cfg := config.Load()
	log.Println("Configuration loaded successfully")

	pool := database.NewConn(cfg.DBURL)
	defer pool.Close()

	store := postgres.NewStore(pool)
	log.Println("Repository layer initialized")

	teamService := service.NewTeamService(store)
	userService := service.NewUserService(store)
	prService := service.NewPRService(store)
	log.Println("UseCase layer initialized")

	handler := httpDelivery.NewHandler(teamService, userService, prService)

	e := httpDelivery.NewRouter(handler)
	log.Println("HTTP handlers initialized")

	port := ":" + cfg.Port
	go func() {
		log.Printf("Starting server on %s", port)
		if err := e.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
