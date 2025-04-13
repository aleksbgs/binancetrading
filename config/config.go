package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Symbols []string `yaml:"symbols"`
	Binance struct {
		RetryDelay time.Duration `yaml:"retry_delay"`
		MaxRetries int           `yaml:"max_retries"`
	} `yaml:"binance"`
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	APIPort string
}

// Load loads the configuration from a YAML file and environment variables
func Load(filePath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Load database configuration from environment variables
	cfg.Database.Host = os.Getenv("DATABASE_HOST")
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}

	cfg.Database.Port = os.Getenv("DATABASE_PORT")
	if cfg.Database.Port == "" {
		cfg.Database.Port = "5432"
	}

	cfg.Database.User = os.Getenv("DATABASE_USER")
	if cfg.Database.User == "" {
		cfg.Database.User = "postgres"
	}

	cfg.Database.Password = os.Getenv("DATABASE_PASSWORD")
	if cfg.Database.Password == "" {
		cfg.Database.Password = "postgres"
	}

	cfg.Database.Name = os.Getenv("DATABASE_NAME")
	if cfg.Database.Name == "" {
		cfg.Database.Name = "trading"
	}

	cfg.Database.SSLMode = os.Getenv("DATABASE_SSLMODE")
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable"
	}

	cfg.APIPort = os.Getenv("API_PORT")
	if cfg.APIPort == "" {
		cfg.APIPort = "8080"
	}

	return &cfg, nil
}

// GetDSN returns the PostgreSQL DSN string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode,
	)
}
