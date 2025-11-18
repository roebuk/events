package main

import (
	"context"
	"firecrest-go/tutorial"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	logger *slog.Logger
	db     *tutorial.Queries
}

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	dbpool, dbErr := pgxpool.New(context.Background(), "postgres://postgres:postgres@127.0.0.1:5432/firecrest")

	if dbErr != nil {
		logger.Error(dbErr.Error())
		os.Exit(1)
	}

	defer dbpool.Close()

	queries := tutorial.New(dbpool)

	app := &application{
		logger: logger,
		db:     queries,
	}

	authors, err := queries.ListAuthors(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}

	log.Println(authors)

	for _, author := range authors {
		fmt.Printf("Author: %d, %s, %s\n", author.ID, author.Name, author.Bio.String)
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
	err = srv.ListenAndServe()

	fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	os.Exit(1)
}
