package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"gobooking/internal/config"
	"gobooking/internal/http/router"
	"gobooking/internal/repository/redis"
	"gobooking/internal/worker"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run database migrations
	if err := runMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Connect to PostgreSQL
	db, err := connectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redisClient := redis.NewRedisClient(
		fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	defer redisClient.Close()

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Create router
	router := router.NewRouter(cfg, db, redisClient)

	// Start background worker
	workerService := worker.NewWorkerService(
		worker.NewHoldReleaseWorker(db, cfg.Worker.HoldReleaseInterval),
		worker.NewBookingExpireWorker(db, cfg.Worker.BookingExpireInterval),
	)

	if err := workerService.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router.Engine(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	if err := workerService.Stop(); err != nil {
		log.Printf("Worker shutdown completed with error: %v", err)
	}

	log.Println("Server exited")
}

func connectDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := cfg.Postgres.DSN()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.Postgres.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.Postgres.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.Postgres.ConnMaxLifetime

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(cfg *config.Config) error {
	dsn := cfg.Postgres.DSN()
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
