#!/bin/sh
set -e

if [ -d "db/migrations" ] && [ -n "$DATABASE_URL" ]; then
  echo "[INIT] Running database migrations..."
  for f in db/migrations/*.up.sql; do
    echo "  -> $(basename "$f")"
    psql "$DATABASE_URL" -f "$f" 2>&1 || true
  done
  echo "[INIT] Migrations complete."
fi

echo "[INIT] Starting API server..."
exec ./alizia-inclusion-api
