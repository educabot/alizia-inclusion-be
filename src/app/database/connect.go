// Package database centraliza la conexión a Postgres (con reintento/backoff y
// pooling) para reusarla entre los entrypoints: el server web (cmd) y los jobs
// batch (cmd/summarizer).
package database

import (
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/educabot/alizia-inclusion-be/config"
)

// connectMaxAttempts y connectBackoff acotan el reintento de conexión al boot.
// Las DB gestionadas detrás de un proxy (ej. Railway) suelen resetear la primera
// conexión en frío; un único intento haría fallar el arranque.
const (
	connectMaxAttempts = 10
	connectBackoff     = 2 * time.Second
)

// Connect abre la conexión, la verifica con Ping (reintentando con backoff) y
// aplica el pooling de config. Devuelve error en vez de matar el proceso para
// que cada entrypoint decida cómo manejar el fallo.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var lastErr error

	for attempt := 1; attempt <= connectMaxAttempts; attempt++ {
		db, lastErr = openAndPing(cfg)
		if lastErr == nil {
			return db, nil
		}
		slog.Warn("database connection attempt failed",
			"attempt", attempt, "max_attempts", connectMaxAttempts, "error", lastErr)
		if attempt < connectMaxAttempts {
			time.Sleep(connectBackoff)
		}
	}
	return nil, lastErr
}

func openAndPing(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	return db, nil
}
