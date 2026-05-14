#!/bin/bash
set -e

CONTAINER="alizia-inclusion-postgres"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "Seeding database..."
docker exec -i $CONTAINER psql -U postgres -d alizia_inclusion < "$PROJECT_DIR/db/seeds/seed.sql"
echo "Seeding complete."
