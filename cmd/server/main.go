package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"medods/internal/config"
	"medods/internal/http-server/handlers"
	"medods/internal/http-server/routes"
	kafka "medods/internal/infrustructure/kafka/producer"
	"medods/internal/infrustructure/logger/handlers/slogpretty"
	"medods/internal/infrustructure/logger/sl"
	"medods/internal/repositories/postgres"
	"medods/internal/services"
	"medods/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "medods/docs"

	swaggerFiles "github.com/swaggo/files"
)

// go run ./cmd/server/main.go -config=./config/config.yaml -migration=file://./db/migrations
func main() {
	// read flags
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	migrationPath := flag.String("migration", "file://db/migrations", "Path to the migration directory")
	flag.Parse()

	// read configuration fila
	cfg := config.MustLoadConfig(*configPath)

	// create logger
	logger := setupPrettyLog()

	// create database connection
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	dbConn, err := pgx.Connect(context.Background(), dbConnStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL: %v", err))
	}

	// make migration
	postgres.MustRunMigration(dbConnStr, *migrationPath)

	// create a copy of repository copy
	userRepo := postgres.NewUserRepository(dbConn)

	// create Kafka producer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := kafka.NewKafkaProducer([]string{cfg.Kafka.Providers}, config)
	if err != nil {
		logger.Error("failed create kafka producer", sl.Err(err))
		panic("failed create kafka producer")
	}

	//create a copy of jwt utils
	jwtUtils := utils.NewJWTUtils(cfg.JWTSecret)
	// create a auth services copy
	services := services.NewAuthService(jwtUtils, userRepo, producer)

	// create a handlers copy
	handlers := handlers.NewAuthHandlers(logger, services)

	// create router copy
	router := routes.NewRouter()
	router.SetupRouter(handlers)

	// set up swagger
	router.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := http.Server{
		Addr:         cfg.Port,
		Handler:      router.Engine,
		WriteTimeout: time.Duration(cfg.Timeouts.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Timeouts.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Timeouts.IdleTimeout) * time.Second,
	}

	// realize a gracefull shutdown
	errChan := make(chan error, 1)

	go func() {
		logger.Info("Starting server on", "port", cfg.Port)
		errChan <- server.ListenAndServe()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigint:
		logger.Debug("Caught signal", "signal", sig)
	case err := <-errChan:
		logger.Error("error listen and serve", sl.Err(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", sl.Err(err))
	}

	logger.Info("Server stopped gracefully")

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
