package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subServ/internal/config"
	"subServ/internal/repository/postgres"
	"subServ/internal/service"
	transporthttp "subServ/internal/transport/http"
	"subServ/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 1. конфиг
	cfg := config.MustLoad()

	// 2. логгер
	l, err := logger.New(cfg.App.LogLevel)
	if err != nil {
		log.Fatalf("logger error: %s", err)
	}

	l.Info("starting application", "env", cfg.App.Env)

	// 3. база данных
	db, err := postgres.NewPool(cfg.Database)
	if err != nil {
		l.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	l.Info("connected to database")

	// 4. миграции
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		l.Error("failed to init migrations", "err", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Error("failed to run migrations", "err", err)
		os.Exit(1)
	}

	l.Info("migrations applied")

	// 5. слои приложения
	repo := postgres.NewSubscriptionRepository(db)
	svc := service.NewSubscriptionService(repo, l) // ← передаём логгер
	handler := transporthttp.NewSubscriptionHandler(svc, l)

	// 6. роутер
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// 7. сервер с таймаутами + middleware
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      transporthttp.LoggingMiddleware(l, mux), // ← оборачиваем
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 8. graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		l.Info("server started", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	l.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		l.Error("forced shutdown", "err", err)
		os.Exit(1)
	}

	l.Info("server stopped")
}
