package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/johnroshan2255/auth-service/internal/config"
	"github.com/johnroshan2255/auth-service/internal/repository"
	"github.com/johnroshan2255/auth-service/internal/service"
	"github.com/johnroshan2255/auth-service/internal/transport/grpc"

	authv1 "github.com/johnroshan2255/auth-service/internal/transport/grpc/auth/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	dbpool, err := pgxpool.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer dbpool.Close()

	userRepo := repository.NewPostgresUserRepo(dbpool)
	authService := service.NewAuthService(userRepo)
	handler := grpc.NewAuthHandler(authService)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(authUnaryInterceptor),
	)
	authv1.RegisterAuthServiceServer(server, handler)

	port := cfg.Port
	if port == "" {
		port = ":1000"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Auth Service running on %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// authUnaryInterceptor is middleware for JWT validation
func authUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// skip login endpoint
	if info.FullMethod == "/auth.v1.AuthService/Login" {
		return handler(ctx, req)
	}

	// validate JWT from metadata if needed here
	return handler(ctx, req)
}
