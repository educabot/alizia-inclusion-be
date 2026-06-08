#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Migraciones con golang-migrate (tracking de versión en schema_migrations).
# Instalar el CLI una vez:  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#
# DB_URL: si no está seteada, usa la DB local de docker-compose.
#   Local:   postgres://postgres:postgres@localhost:5481/alizia_inclusion?sslmode=disable
#   Railway: exportar DB_URL="postgres://...@<host>:<port>/railway?sslmode=require"
DB_URL="${DB_URL:-postgres://postgres:postgres@localhost:5481/alizia_inclusion?sslmode=disable}"

if ! command -v migrate >/dev/null 2>&1; then
    echo "ERROR: 'migrate' no está instalado o no está en el PATH."
    echo "Instalalo con: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    echo "(asegurate de tener \$(go env GOPATH)/bin en el PATH)"
    exit 1
fi

echo "Running migrations (golang-migrate) against ${DB_URL%%:*}://...@${DB_URL##*@}"
migrate -path "$PROJECT_DIR/db/migrations" -database "$DB_URL" up
echo "Migrations complete."
