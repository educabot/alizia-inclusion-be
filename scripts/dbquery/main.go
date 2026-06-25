// Runner de SELECT ad-hoc contra una DB remota (ej. Railway) sin psql.
// Uso: DB_URL="postgres://..." go run ./scripts/dbquery "SELECT ..."
package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	dsn := os.Getenv("DB_URL")
	if dsn == "" || len(os.Args) < 2 {
		return fmt.Errorf("uso: DB_URL=... go run ./scripts/dbquery \"SELECT ...\"")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer func() { _ = db.Close() }()
	db.SetMaxOpenConns(1)

	var rows *sql.Rows
	for i := 0; i < 30; i++ {
		rows, err = db.Query(os.Args[1])
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	cols, _ := rows.Columns()
	fmt.Println(strings.Join(cols, " | "))
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		out := make([]string, len(cols))
		for i, v := range vals {
			out[i] = fmt.Sprintf("%v", v)
		}
		fmt.Println(strings.Join(out, " | "))
	}
	return rows.Err()
}
