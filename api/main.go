package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"topikgram/api/internal/migrations"
	"topikgram/api/internal/store"
	"topikgram/api/internal/transport"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	transport.InitJWTSecret()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	prefix := os.Getenv("API_PREFIX")
	if prefix == "" {
		prefix = "api/v1"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := store.NewPool(ctx, connStr)
	if err != nil {
		log.Fatalf("Database failed: %v", err)
	}
	defer db.Close()

	m := migrations.New(db)
	if err := m.Up(ctx); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	mux := http.NewServeMux()

	transport.NewAuthHandler(mux, db).RegisterRoutes(prefix)
	transport.NewUserHandler(mux, db).RegisterRoutes(prefix)
	transport.NewPostHandler(mux, db).RegisterRoutes(prefix)
	transport.NewCommentHandler(mux, db).RegisterRoutes(prefix)

	mux.HandleFunc("GET /health", transport.Health)

	handler := transport.CorsMiddleware(
		transport.BodyLimitMiddleware(1_048_576)(
			transport.RequestTimeoutMiddleware(30*time.Second)(
				transport.LogMiddleware(mux),
			),
		),
	)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("server shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited gracefully")
}
