package main

import (
	"context"
	"database/sql"
	"design-url-shortner/internal/config"
	"design-url-shortner/internal/handler"
	"design-url-shortner/internal/repository/postgres"
	"design-url-shortner/internal/repository/redis"
	"design-url-shortner/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	redisclient "github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load the config: %v", err)
	}

	// Initialise PostgreSQL
	db, err := sql.Open("postgres", cfg.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Initialize database tables
	if err := postgres.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	rdb := redisclient.NewClient(&redisclient.Options{
		Addr: cfg.GetRedisAddr(),
	})
	defer rdb.Close()

	// Test Redis Connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialise repository
	postgresRepo := postgres.NewURLRepository(db)
	redisRepo := redis.NewURLRepository(rdb)

	// Initialize service
	urlService := service.NewURLService(postgresRepo, redisRepo, cfg.BaseURL)

	// Initialize handler
	urlHandler := handler.NewURLHandler(urlService)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/urls", urlHandler.CreateURL)
		r.Get("/urls/{shortCode}/stats", urlHandler.GetURLStats)
	})

	// Redirect route

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	// Graceful shutdown
	log.Printf("Shutting down server....")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server exiting")

}
