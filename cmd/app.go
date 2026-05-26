package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/educabot/alizia-inclusion-be/config"
	appweb "github.com/educabot/alizia-inclusion-be/src/app/web"
	appmw "github.com/educabot/alizia-inclusion-be/src/entrypoints/middleware"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	initLogger(cfg)

	db := connectDB(cfg)

	repos := NewRepositories(db, cfg)
	usecases := NewUseCases(repos)
	handlers := NewHandlers(usecases, cfg)

	engine := gin.Default()
	engine.Use(appmw.RequestLogger())
	engine.Use(cors.New(cors.Config{
		AllowOriginFunc:  buildOriginChecker(cfg.AllowedOrigins),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	engine.GET("/health", healthHandler(db))

	appweb.ConfigureMappings(engine, handlers, cfg)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	return &App{
		cfg: cfg,
		db:  db,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: engine,
		},
	}
}

func (a *App) Run() {
	go func() {
		slog.Info("server listening", "addr", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}

func (a *App) Close() {
	if sqlDB, err := a.db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

// buildOriginChecker returns a function that checks if an origin is allowed.
// Supports exact matches and wildcard subdomain patterns like "https://*.example.com".
func buildOriginChecker(allowed []string) func(string) bool {
	for _, o := range allowed {
		if o == "*" {
			return func(string) bool { return true }
		}
	}
	return func(origin string) bool {
		for _, o := range allowed {
			if o == origin {
				return true
			}
			if strings.Contains(o, "*") {
				suffix := strings.Replace(o, "*", "", 1)
				scheme := ""
				if idx := strings.Index(suffix, "://"); idx != -1 {
					scheme = suffix[:idx+3]
					suffix = suffix[idx+3:]
				}
				if strings.HasPrefix(origin, scheme) && strings.HasSuffix(origin, suffix) {
					return true
				}
			}
		}
		return false
	}
}

func initLogger(cfg *config.Config) {
	var handler slog.Handler
	if cfg.Env == "local" || cfg.Env == "test" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	slog.SetDefault(slog.New(handler))
}

func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "degraded",
				"db":     "error",
			})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "degraded",
				"db":     "unreachable",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     "ok",
		})
	}
}

func connectDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get sql.DB", "error", err)
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	slog.Info("database connected")
	return db
}
