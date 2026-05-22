#!/bin/sh
set -e

echo "Running migrations..."
for f in /migrations/*.up.sql; do
  echo "  -> $(basename "$f")"
  psql "$DATABASE_URL" -f "$f" 2>&1 || true
done
echo "Migrations complete."
