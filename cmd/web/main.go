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
	"firecrest/internal/config"
	"firecrest/internal/repository"
	"firecrest/internal/service"
)

type application struct {
	config         *config.Config
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
	// Load environment-specific .env file
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Try to load environment-specific .env file first, then fall back to .env
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found, using system environment variables")
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logger.Info("Starting application", "environment", cfg.Environment)

	// Build connection string from config
	dsn := cfg.DatabaseDSN()

	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer dbpool.Close()

	queries := db.New(dbpool)

	// Initialize session manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(dbpool)
	sessionManager.Lifetime = time.Duration(cfg.Session.LifetimeHrs) * time.Hour
	sessionManager.Cookie.Name = cfg.Session.CookieName
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = cfg.Session.SecureCookie

	// Initialize repositories
	eventRepo := repository.NewEventRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	// Initialize services
	eventService := service.NewEventService(eventRepo)
	userService := service.NewUserService(userRepo)

	app := &application{
		config:         cfg,
		logger:         logger,
		sessionManager: sessionManager,
		eventService:   eventService,
		userService:    userService,
		authService:    authService,
	}

	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        app.routes(),
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(cfg.Server.IdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	logger.Info("Starting server", "port", cfg.Server.Port, "environment", cfg.Environment)
	return srv.ListenAndServe()
}
