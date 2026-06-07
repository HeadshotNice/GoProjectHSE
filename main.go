package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"HSE/internal/config"
	"HSE/internal/queue/rabbitmq"
	"HSE/internal/repository/postgres"
	"HSE/internal/transport/httpapi"
	"HSE/internal/usecase"
	"HSE/internal/worker"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.FromEnv()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := postgres.Open(ctx, cfg.DB)
	if err != nil {
		logger.Fatalf("db open: %v", err)
	}
	defer db.Close()

	if err := postgres.InitSchema(ctx, db); err != nil {
		logger.Fatalf("db init schema: %v", err)
	}

	repos := postgres.NewRepositories(db)

	var events *rabbitmq.Client
	if cfg.RabbitMQ.Enabled {
		events, err = rabbitmq.Open(cfg.RabbitMQ)
		if err != nil {
			logger.Fatalf("rabbitmq open: %v", err)
		}
		defer events.Close()
	}

	uc := usecase.New(
		repos.Test,
		repos.DBTest,
		repos.Users,
		repos.Orders,
		repos.Docs,
		events,
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTIssuer,
		cfg.Auth.JWTTTL,
	)

	if events != nil {
		documentWorker := worker.NewDocumentReviewWorker(events, uc, logger)
		go func() {
			if err := documentWorker.Run(ctx); err != nil && err != context.Canceled {
				logger.Printf("document review worker stopped: %v", err)
			}
		}()
	}

	handler := httpapi.NewHandler(uc, cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer)

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Printf("listening on %s", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Printf("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("shutdown error: %v", err)
	} else {
		logger.Printf("shutdown complete")
	}
}
