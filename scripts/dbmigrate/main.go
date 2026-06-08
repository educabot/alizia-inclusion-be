// Runner de migraciones/seeds contra una DB remota (ej. Railway) sin psql.
// Aplica archivos .sql idempotentes en orden, tolerando el proxy inestable.
//
// Uso:
//   DB_URL="postgres://..." go run ./scripts/dbmigrate <archivo.sql> [archivo2.sql ...]
//
// Cada archivo se ejecuta entero (multi-statement, simple query protocol de lib/pq).
// Si la conexión se resetea, reintenta el archivo completo (los .sql son idempotentes).
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		fmt.Println("falta DB_URL")
		os.Exit(1)
	}
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Println("uso: go run ./scripts/dbmigrate <archivo.sql> [...]")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("open:", err)
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	// warmup: la primera query en frío suele resetear; reintento hasta tener conexión viva.
	for i := 0; i < 30; i++ {
		var x int
		if db.QueryRow(`SELECT 1`).Scan(&x) == nil {
			break
		}
	}

	failed := false
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Printf("FAIL  %s  (no se pudo leer: %v)\n", f, err)
			failed = true
			continue
		}
		var lastErr error
		ok := false
		for attempt := 0; attempt < 25; attempt++ {
			if _, err := db.Exec(string(content)); err == nil {
				ok = true
				break
			} else {
				lastErr = err
			}
		}
		if ok {
			fmt.Printf("OK    %s\n", f)
		} else {
			fmt.Printf("FAIL  %s  -> %v\n", f, lastErr)
			failed = true
		}
	}

	if failed {
		os.Exit(2)
	}
	fmt.Println("Listo.")
}
