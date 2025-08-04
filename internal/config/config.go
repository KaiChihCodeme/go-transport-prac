package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `envconfig:"SERVER"`
	
	// Database configuration
	Database DatabaseConfig `envconfig:"DATABASE"`
	
	// Redis configuration
	Redis RedisConfig `envconfig:"REDIS"`
	
	// MinIO configuration
	MinIO MinIOConfig `envconfig:"MINIO"`
	
	// Logging configuration
	Logging LoggingConfig `envconfig:"LOGGING"`
	
	// Development configuration
	Development DevelopmentConfig `envconfig:"DEV"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	HTTPPort     int           `envconfig:"HTTP_PORT" default:"8080"`
	GRPCPort     int           `envconfig:"GRPC_PORT" default:"8081"`
	WSPort       int           `envconfig:"WS_PORT" default:"8082"`
	GraphQLPort  int           `envconfig:"GRAPHQL_PORT" default:"9090"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"30s"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT" default:"120s"`
	Host         string        `envconfig:"HOST" default:"localhost"`
	TLSEnabled   bool          `envconfig:"TLS_ENABLED" default:"false"`
	CertFile     string        `envconfig:"CERT_FILE"`
	KeyFile      string        `envconfig:"KEY_FILE"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string        `envconfig:"HOST" default:"localhost"`
	Port         int           `envconfig:"PORT" default:"5432"`
	Username     string        `envconfig:"USERNAME" default:"transport_user"`
	Password     string        `envconfig:"PASSWORD" default:"transport_pass"`
	Name         string        `envconfig:"NAME" default:"transport_db"`
	SSLMode      string        `envconfig:"SSL_MODE" default:"disable"`
	MaxOpenConns int           `envconfig:"MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns int           `envconfig:"MAX_IDLE_CONNS" default:"5"`
	MaxLifetime  time.Duration `envconfig:"MAX_LIFETIME" default:"300s"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string        `envconfig:"HOST" default:"localhost"`
	Port         int           `envconfig:"PORT" default:"6379"`
	Password     string        `envconfig:"PASSWORD"`
	Database     int           `envconfig:"DATABASE" default:"0"`
	MaxRetries   int           `envconfig:"MAX_RETRIES" default:"3"`
	PoolSize     int           `envconfig:"POOL_SIZE" default:"10"`
	DialTimeout  time.Duration `envconfig:"DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"3s"`
}

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint        string `envconfig:"ENDPOINT" default:"localhost:9000"`
	AccessKeyID     string `envconfig:"ACCESS_KEY_ID" default:"minioadmin"`
	SecretAccessKey string `envconfig:"SECRET_ACCESS_KEY" default:"minioadmin"`
	UseSSL          bool   `envconfig:"USE_SSL" default:"false"`
	BucketName      string `envconfig:"BUCKET_NAME" default:"transport-data"`
	Region          string `envconfig:"REGION" default:"us-east-1"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level       string `envconfig:"LEVEL" default:"info"`
	Format      string `envconfig:"FORMAT" default:"json"`
	OutputPaths string `envconfig:"OUTPUT_PATHS" default:"stdout"`
	Development bool   `envconfig:"DEVELOPMENT" default:"false"`
}

// DevelopmentConfig holds development-specific configuration
type DevelopmentConfig struct {
	Enabled         bool `envconfig:"ENABLED" default:"false"`
	MockServices    bool `envconfig:"MOCK_SERVICES" default:"false"`
	EnableProfiling bool `envconfig:"ENABLE_PROFILING" default:"false"`
	EnableMetrics   bool `envconfig:"ENABLE_METRICS" default:"true"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	
	// Process environment variables
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server ports
	if c.Server.HTTPPort <= 0 || c.Server.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", c.Server.HTTPPort)
	}
	
	if c.Server.GRPCPort <= 0 || c.Server.GRPCPort > 65535 {
		return fmt.Errorf("invalid gRPC port: %d", c.Server.GRPCPort)
	}
	
	// Validate database configuration
	if c.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	
	if c.Database.Username == "" {
		return fmt.Errorf("database username cannot be empty")
	}
	
	// Validate logging level
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	if !contains(validLevels, strings.ToLower(c.Logging.Level)) {
		return fmt.Errorf("invalid logging level: %s", c.Logging.Level)
	}
	
	// Validate TLS configuration
	if c.Server.TLSEnabled {
		if c.Server.CertFile == "" || c.Server.KeyFile == "" {
			return fmt.Errorf("TLS cert file and key file must be specified when TLS is enabled")
		}
		
		// Check if files exist
		if _, err := os.Stat(c.Server.CertFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS cert file does not exist: %s", c.Server.CertFile)
		}
		
		if _, err := os.Stat(c.Server.KeyFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS key file does not exist: %s", c.Server.KeyFile)
		}
	}
	
	return nil
}

// DatabaseURL returns the database connection URL
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// RedisAddr returns the Redis address
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// MinIOEndpoint returns the MinIO endpoint
func (c *Config) MinIOEndpoint() string {
	return c.MinIO.Endpoint
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Development.Enabled
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets an environment variable as integer with a default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool gets an environment variable as boolean with a default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// GetEnvDuration gets an environment variable as duration with a default value
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}