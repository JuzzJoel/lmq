package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/user/lmq/backend/config"
	"github.com/user/lmq/backend/database"
	"github.com/user/lmq/backend/handlers"
	appmw "github.com/user/lmq/backend/middleware"
	"github.com/user/lmq/backend/services"
	"github.com/user/lmq/backend/spa"
)

func main() {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Main]: Failed to load config: %v", err)
	}

	// 2. Initialize Postgres pool + run migrations
	ctx := context.Background()
	isProd := os.Getenv("APP_ENV") == "production" || os.Getenv("APP_ENV") == "prod" || os.Getenv("ENV") == "production" || os.Getenv("ENV") == "prod" || os.Getenv("RENDER") != ""

	pool, err := database.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		if isProd {
			log.Fatalf("[Main]: Database init failed: %v", err)
		}
		log.Printf("[EMERGENCY BYPASS] Database/Cache unavailable. Running in decoupled headless testing mode on port %s. (DB Error: %v)", cfg.Port, err)
	}
	if pool != nil {
		defer pool.Close()
	}

	if pool != nil {
		if err := database.RunMigrations(ctx, pool); err != nil {
			if isProd {
				log.Fatalf("[Main]: Migrations failed: %v", err)
			}
			log.Printf("[Main]: Migrations failed: %v", err)
		}
	}

	// 3. Initialize Redis client
	rdb, err := database.NewRedisClient(cfg.RedisURL)
	if err != nil {
		if isProd {
			log.Fatalf("[Main]: Redis init failed: %v", err)
		}
		log.Printf("[EMERGENCY BYPASS] Redis unavailable. (Error: %v)", err)
	}
	if rdb != nil {
		defer rdb.Close()
	}

	// 4. Initialize GeoIP service
	geoIPService := services.NewGeoIPService(cfg.GeoIPDBPath)
	defer geoIPService.Close()

	// 5. Initialize services (analytics)
	analyticsService := services.NewAnalyticsService(pool, geoIPService)

	// 6. Initialize handlers
	shortenHandler := handlers.NewShortenHandler(pool, rdb)
	csvShortenHandler := handlers.NewCsvShortenHandler(pool, rdb)
	redirectHandler := handlers.NewRedirectHandler(pool, rdb, analyticsService)
	analyticsHandler := handlers.NewAnalyticsHandler(pool)
	verifyHandler := handlers.NewVerifyHandler(pool)

	// 7. Build Chi router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(appmw.StrictCORSMiddleware())

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs", handlers.HandleGetDocs)
		r.With(appmw.LocalRateLimiter()).Post("/shorten", shortenHandler.HandleShorten)
		r.With(appmw.LocalRateLimiter()).Post("/shorten/csv", csvShortenHandler.HandleCsvShorten)
		r.Post("/verify-password", verifyHandler.HandleVerify)
		
		r.Route("/analytics", func(r chi.Router) {
			r.Use(appmw.AdminAuthMiddleware)
			r.Get("/", analyticsHandler.HandleGetAnalytics)
			r.Get("/overview", analyticsHandler.HandleGetOverview)
			r.Get("/links", analyticsHandler.HandleListLinks)
			r.Get("/export", analyticsHandler.HandleExportAnalytics)
		})
	})

	r.Get("/{token}", redirectHandler.HandleRedirect)

	r.NotFound(spa.Handler())

	// 8. Start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("[Main]: Server listening on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Main]: Server error: %v", err)
		}
	}()

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Main]: Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[Main]: Server forced to shutdown: %v", err)
	}
	log.Println("[Main]: Server exiting")
}
