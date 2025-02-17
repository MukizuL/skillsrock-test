// Package web реализует REST API для управления задачами.
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/MukizuL/skillsrock-test/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"os"
)

// application is used for providing necessary dependencies
type application struct {
	logger *slog.Logger
	tasks  *models.TaskModel
}

// Main creates database connection and starts a fiber server
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	addr := getEnv("ADDR", ":8080")
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "admin"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "postgres"))

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	err = Migrate(connString)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger: logger,
		tasks:  &models.TaskModel{DB: db},
	}

	server := app.routes()

	logger.Info("starting server", "addr", addr)

	err = server.Listen(addr)
	logger.Error(err.Error())
	os.Exit(1)
}

// Migrate helps to migrate databases
func Migrate(DSN string) error {
	m, err := migrate.New("file:///migrations", DSN)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// getEnv searches env variable by key. If none found, returns fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Key is empty", "key=", key)
	return fallback
}
