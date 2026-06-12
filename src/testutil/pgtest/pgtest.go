//go:build integration

// Package pgtest provides a real PostgreSQL instance (via testcontainers) for
// repository integration tests. It starts one container per test process, applies
// the real db/migrations/*.up.sql, and hands out transaction-scoped *gorm.DB
// handles that roll back after each test for isolation.
//
// Tests using this package must carry the `integration` build tag and require a
// running Docker daemon. The image is pgvector/pgvector:pg16 because migration
// 000021 enables the `vector` extension, absent from the stock postgres image.
package pgtest

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/educabot/alizia-inclusion-be/src/testutil"
)

// OtherOrgID is a second seeded organization used to assert tenant isolation:
// a repo scoped to testutil.TestOrgID must never read OtherOrgID's rows.
var OtherOrgID = uuid.MustParse("b2c3d4e5-f6a7-8901-bcde-f23456789012")

var (
	once    sync.Once
	sharedDB *gorm.DB
	setupErr error
)

// DB lazily starts the shared Postgres container (once per test process), applies
// every migration, and returns a GORM handle bound to that database. It fails the
// test if Docker is unavailable or migrations error.
func DB(t testing.TB) *gorm.DB {
	t.Helper()
	once.Do(start)
	if setupErr != nil {
		t.Fatalf("pgtest: %v", setupErr)
	}
	return sharedDB
}

// Tx returns a transaction-scoped *gorm.DB and registers a rollback on cleanup, so
// each test sees a clean database with no cross-test state. Pass the returned handle
// to the repository under test.
func Tx(t *testing.T) *gorm.DB {
	t.Helper()
	db := DB(t)
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("pgtest: begin tx: %v", tx.Error)
	}
	t.Cleanup(func() { tx.Rollback() })
	return tx
}

func start() {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "pgvector/pgvector:pg16",
		tcpostgres.WithDatabase("alizia_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		setupErr = err
		return
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		setupErr = err
		return
	}

	if err := applyMigrations(dsn); err != nil {
		setupErr = err
		return
	}

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		setupErr = err
		return
	}

	// Seed the tenant rows once (committed, outside any test tx) so every test can
	// reference them via FK. Per-test data is created inside the rolled-back tx.
	if err := gdb.Exec(
		"INSERT INTO organizations (id, name) VALUES (?, 'Test Org'), (?, 'Other Org') ON CONFLICT (id) DO NOTHING",
		testutil.TestOrgID, OtherOrgID,
	).Error; err != nil {
		setupErr = err
		return
	}
	sharedDB = gdb
}

// applyMigrations runs every db/migrations/*.up.sql in numeric order using lib/pq's
// simple query protocol, which (unlike the extended protocol) accepts the
// multi-statement files as-is — mirroring scripts/dbmigrate.
func applyMigrations(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	files, err := upMigrations()
	if err != nil {
		return err
	}
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(content)); err != nil {
			return &migrationError{file: filepath.Base(f), err: err}
		}
	}
	return nil
}

// upMigrations returns the absolute paths of the *.up.sql migration files, sorted by
// their numeric prefix, located by walking up from this source file to the repo root.
func upMigrations() ([]string, error) {
	root, err := repoRoot()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(root, "db", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(files) // zero-padded numeric prefixes sort lexicographically
	return files, nil
}

// repoRoot walks up from this file's directory until it finds go.mod.
func repoRoot() (string, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

type migrationError struct {
	file string
	err  error
}

func (e *migrationError) Error() string { return "migration " + e.file + ": " + e.err.Error() }
func (e *migrationError) Unwrap() error  { return e.err }
