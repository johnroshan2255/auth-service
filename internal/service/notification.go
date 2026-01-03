package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/johnroshan2255/auth-service/internal/config"
	notificationv1 "github.com/johnroshan2255/core-service/proto/notification/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type CoreNotificationClient struct {
	conn       *grpc.ClientConn
	addr       string
	serviceKey string
}

func NewCoreNotificationClient(addr, serviceKey string, useTLS bool) (*CoreNotificationClient, error) {
	if addr == "" {
		return nil, fmt.Errorf("core notification service address is required")
	}

	var creds credentials.TransportCredentials
	if useTLS {
		config := &tls.Config{
			InsecureSkipVerify: false,
		}
		creds = credentials.NewTLS(config)
		log.Printf("Connecting to core notification service with TLS at %s", addr)
	} else {
		creds = insecure.NewCredentials()
		log.Printf("WARNING: Connecting to core notification service without TLS at %s", addr)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to core notification service: %w", err)
	}

	client := &CoreNotificationClient{
		conn:       conn,
		addr:       addr,
		serviceKey: serviceKey,
	}

	log.Printf("gRPC client connected to core notification service at %s", addr)
	return client, nil
}

// Close closes the gRPC connection
func (c *CoreNotificationClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *CoreNotificationClient) createContextWithAuth(ctx context.Context) context.Context {
	if c.serviceKey == "" {
		return ctx
	}
	md := metadata.New(map[string]string{
		"service-key": c.serviceKey,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *CoreNotificationClient) NotifyUserCreated(ctx context.Context, userUUID, email, username string) error {
	ctx = c.createContextWithAuth(ctx)
	
	client := notificationv1.NewNotificationServiceClient(c.conn)
	req := &notificationv1.UserCreatedRequest{
		UserUuid: userUUID,
		Email:    email,
		Username: username,
	}
	
	_, err := client.NotifyUserCreated(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to notify user creation: %w", err)
	}
	
	log.Printf("Successfully notified core service: User created - UUID: %s, Email: %s, Username: %s", userUUID, email, username)
	return nil
}

type NotificationService struct {
	client *CoreNotificationClient
}

func NewNotificationService(cfg *config.Config) (*NotificationService, error) {
	if cfg.CoreNotificationServiceAddr == "" {
		log.Println("Warning: CORE_NOTIFICATION_SERVICE_ADDR not set. Core notification calls will be disabled.")
		return nil, nil
	}

	client, err := NewCoreNotificationClient(cfg.CoreNotificationServiceAddr, cfg.ServiceKey, cfg.TLSEnabled)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize notification service: %w", err)
	}

	return &NotificationService{
		client: client,
	}, nil
}

func (ns *NotificationService) SetupAuthService(authService *AuthService) {
	if ns != nil && ns.client != nil {
		authService.SetCoreNotificationClient(ns.client)
	}
}

func (ns *NotificationService) Close() error {
	if ns != nil && ns.client != nil {
		return ns.client.Close()
	}
	return nil
}

func (ns *NotificationService) IsEnabled() bool {
	return ns != nil && ns.client != nil
}

