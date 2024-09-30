package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"medods/internal/config"
	"medods/internal/logger/handlers/slogpretty"
	"medods/internal/repositories/postgres"
	"os"

	"github.com/jackc/pgx/v5"
)

// go run ./cmd/server/main.go -config=./config/config.yaml -migration=file://./db/migrations
func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	migrationPath := flag.String("migration", "file://db/migrations", "Path to the migration directory")
	flag.Parse()

	cfg := config.MustLoadConfig(*configPath)

	logger := setupPrettyLog()

	logger.Info("Starting server on", "port", *migrationPath)

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	dbConn, err := pgx.Connect(context.Background(), dbConnStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL: %v", err))
	}

	postgres.MustRunMigration(dbConnStr, *migrationPath)

	userRepo := postgres.NewUserRepository(dbConn)
}

func setupPrettyLog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
