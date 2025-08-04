package wire

import (
	"github.com/google/wire"

	"go-transport-prac/internal/config"
	"go-transport-prac/internal/logger"
)

// ConfigProviderSet provides configuration-related dependencies
var ConfigProviderSet = wire.NewSet(
	ProvideConfig,
)

// LoggerProviderSet provides logging-related dependencies
var LoggerProviderSet = wire.NewSet(
	ProvideLogger,
)

// InfrastructureProviderSet provides all infrastructure dependencies
var InfrastructureProviderSet = wire.NewSet(
	ConfigProviderSet,
	LoggerProviderSet,
)

// SDLProviderSet will contain Schema Definition Language providers
var SDLProviderSet = wire.NewSet(
	// Will be populated when SDL packages are implemented
)

// WebProtocolProviderSet will contain Web Protocol providers
var WebProtocolProviderSet = wire.NewSet(
	// Will be populated when web protocol packages are implemented
)

// ServiceProviderSet combines all service providers
var ServiceProviderSet = wire.NewSet(
	InfrastructureProviderSet,
	SDLProviderSet,
	WebProtocolProviderSet,
)

// TestProviderSet provides dependencies for testing
var TestProviderSet = wire.NewSet(
	ProvideTestConfig,
	ProvideTestLogger,
)

// ProvideTestConfig provides test configuration
func ProvideTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			HTTPPort:    8080,
			GRPCPort:    8081,
			WSPort:      8082,
			GraphQLPort: 9090,
			Host:        "localhost",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test_user",
			Password: "test_pass",
			Name:     "test_db",
			SSLMode:  "disable",
		},
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
		Logging: config.LoggingConfig{
			Level:       "debug",
			Format:      "console",
			OutputPaths: "stdout",
			Development: true,
		},
		Development: config.DevelopmentConfig{
			Enabled: true,
		},
	}
}

// ProvideTestLogger provides test logger
func ProvideTestLogger() (*logger.Logger, error) {
	return logger.NewDevelopment()
}