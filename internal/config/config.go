package config

import (
	"os"
)

type Config struct {
	DBUrl    string
	Port     string
	GRPCPort string
	JWTKey   string
	ServiceKey string
	CoreNotificationServiceAddr string
	// TLS configuration for secure gRPC connections
	TLSCertFile string // Path to TLS certificate file (optional)
	TLSKeyFile  string // Path to TLS key file (optional)
	TLSEnabled  bool   // Enable TLS for gRPC connections
}

func LoadConfig() *Config {
	return &Config{
		DBUrl:     os.Getenv("POSTGRES_URL"),
		Port:      os.Getenv("PORT"),
		GRPCPort:  os.Getenv("GRPC_PORT"),
		JWTKey:    os.Getenv("JWT_KEY"),
		ServiceKey: os.Getenv("SERVICE_KEY"),
		CoreNotificationServiceAddr: os.Getenv("CORE_NOTIFICATION_SERVICE_ADDR"),
		TLSCertFile: os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:  os.Getenv("TLS_KEY_FILE"),
		TLSEnabled:  os.Getenv("TLS_ENABLED") == "true",
	}
}
