package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"firecrest/db"
	"firecrest/internal/repository"
	"firecrest/internal/service"
)

type application struct {
	logger         *slog.Logger
	sessionManager *scs.SessionManager
	eventService   service.EventService
	userService    service.UserService
	authService    service.AuthService
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load .env file in development
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "firecrest")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Build connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	// Initialize session manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Name = "firecrest_session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false // Set to true in production with HTTPS

	// Initialize repositories
	eventRepo := repository.NewEventRepository(queries)
	userRepo := repository.NewUserRepository(queries)
	authRepo := repository.NewAuthRepository(queries)

	// Initialize services
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo)

	app := &application{
		logger:         logger,
		sessionManager: sessionManager,
		eventService:   eventService,
		userService:    userService,
		authService:    authService,
	}

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        app.routes(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	fmt.Println("Running server on :8080")
	return srv.ListenAndServe()
}

// getEnv retrieves the value of an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
