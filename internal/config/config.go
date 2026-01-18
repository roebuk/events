package config

import (
	"fmt"
	"os"
	"strconv"
)

// Environment represents the application environment
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// Config holds all application configuration
type Config struct {
	Environment Environment
	Server      ServerConfig
	Database    DatabaseConfig
	Session     SessionConfig
	CSRF        CSRFConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int // seconds
	WriteTimeout int // seconds
	IdleTimeout  int // seconds
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type SessionConfig struct {
	CookieName   string
	LifetimeHrs  int
	SecureCookie bool
}

type CSRFConfig struct {
	Key            string
	SecureCookie   bool
	TrustedOrigins []string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	env := Environment(getEnv("APP_ENV", "development"))

	config := &Config{
		Environment: env,
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 5),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 120),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "firecrest"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Session: SessionConfig{
			CookieName:   getEnv("SESSION_COOKIE_NAME", "firecrest_session"),
			LifetimeHrs:  getEnvAsInt("SESSION_LIFETIME_HRS", 12),
			SecureCookie: env == Production,
		},
		CSRF: CSRFConfig{
			Key:          os.Getenv("CSRF_KEY"),
			SecureCookie: env != Development,
		},
	}

	// Set trusted origins based on environment
	switch env {
	case Development:
		config.CSRF.TrustedOrigins = []string{"localhost:8080", "127.0.0.1:8080"}
	case Staging:
		config.CSRF.TrustedOrigins = []string{getEnv("STAGING_DOMAIN", "")}
	case Production:
		config.CSRF.TrustedOrigins = []string{getEnv("PRODUCTION_DOMAIN", "")}
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate ensures required configuration values are present
func (c *Config) validate() error {
	if len(c.CSRF.Key) != 32 {
		return fmt.Errorf("CSRF_KEY must be exactly 32 bytes, got %d", len(c.CSRF.Key))
	}

	if c.Environment == Production {
		if !c.CSRF.SecureCookie {
			return fmt.Errorf("CSRF secure cookie must be enabled in production")
		}
		if !c.Session.SecureCookie {
			return fmt.Errorf("session secure cookie must be enabled in production")
		}
		if c.Database.SSLMode == "disable" {
			return fmt.Errorf("database SSL should be enabled in production")
		}
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == Development
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == Production
}

// DatabaseDSN returns the database connection string
func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
