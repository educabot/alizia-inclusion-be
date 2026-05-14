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

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/educabot/alizia-inclusion-be/config"
	appweb "github.com/educabot/alizia-inclusion-be/src/app/web"
)

type App struct {
	cfg    *config.Config
	db     *gorm.DB
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	db := connectDB(cfg)

	repos := NewRepositories(db, cfg)
	usecases := NewUseCases(repos)
	handlers := NewHandlers(usecases, cfg)

	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

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
		log.Printf("[INFO] Server listening on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[INFO] Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Shutdown error: %v", err)
	}
}

func (a *App) Close() {
	if sqlDB, err := a.db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

func connectDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		log.Fatalf("[FATAL] Database connection failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[FATAL] Failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	log.Println("[INFO] Database connected")
	return db
}
