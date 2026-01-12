package main

import (
	"context"
	"firecrest-go/db"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	logger *slog.Logger
	db     *db.Queries
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	dbpool, dbErr := pgxpool.New(context.Background(), "postgres://postgres:postgres@127.0.0.1:5432/firecrest")

	if dbErr != nil {
		logger.Error(dbErr.Error())
		os.Exit(1)
	}

	defer dbpool.Close()

	queries := db.New(dbpool)

	app := &application{
		logger: logger,
		db:     queries,
	}

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        app.routes(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB

	}

	fmt.Println("🚀 Running server on :8080")
	err := srv.ListenAndServe()

	fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	os.Exit(1)
}
