package main

import (
	"log"
	"net"

	"github.com/johnroshan2255/auth-service/internal/config"
	"github.com/johnroshan2255/auth-service/internal/database"
	"github.com/johnroshan2255/auth-service/internal/middleware"
	"github.com/johnroshan2255/auth-service/internal/repository"
	"github.com/johnroshan2255/auth-service/internal/service"
	grpchandler "github.com/johnroshan2255/auth-service/internal/transport/grpc"
	"github.com/johnroshan2255/auth-service/internal/transport/http"
	authv1 "github.com/johnroshan2255/auth-service/proto/auth/v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.LoadConfig()

	if cfg.JWTKey == "" {
		log.Fatal("JWT_KEY environment variable is required")
	}
	service.SetJWTKey(cfg.JWTKey)
	middleware.SetJWTKey(cfg.JWTKey)

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer database.CloseDB(db)

	userRepo := repository.NewPostgresUserRepo(db)
	authService := service.NewAuthService(userRepo)

	// Set service key for backend-to-backend gRPC authentication
	if cfg.ServiceKey != "" {
		middleware.SetServiceKey(cfg.ServiceKey)
	}

	// Initialize notification service
	notificationService, err := service.NewNotificationService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize notification service: %v. Continuing without it.", err)
	} else if notificationService != nil {
		notificationService.SetupAuthService(authService)
		defer notificationService.Close()
	}

	// Start gRPC server in a goroutine
	go func() {
		grpcPort := cfg.GRPCPort
		if grpcPort == "" {
			grpcPort = ":9090"
		}

		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(middleware.BackendAuthInterceptor),
		)

		handler := grpchandler.NewAuthHandler(authService)
		authv1.RegisterAuthServiceServer(grpcServer, handler)

		log.Printf("gRPC server (backend-to-backend) running on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	// Start HTTP server (blocks main thread)
	router := http.SetupRouter(authService)
	port := cfg.Port
	if port == "" {
		port = ":8080"
	}
	log.Printf("HTTP server (frontend) running on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}