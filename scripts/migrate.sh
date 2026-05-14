#!/bin/bash
set -e

CONTAINER="alizia-inclusion-postgres"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "Running migrations..."
for f in "$PROJECT_DIR"/db/migrations/*.up.sql; do
    echo "  Applying $(basename "$f")..."
    docker exec -i $CONTAINER psql -U postgres -d alizia_inclusion < "$f"
done
echo "Migrations complete."
