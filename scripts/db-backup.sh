#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Load environment variables from .env.production if it exists
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"

if [ -f "$PARENT_DIR/.env.production" ]; then
    export $(grep -v '^#' "$PARENT_DIR/.env.production" | xargs)
fi

# Configuration with defaults
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-chatsphere}
CONTAINER_NAME=${CONTAINER_NAME:-chatsphere-postgres-prod}
BACKUP_DIR="${PARENT_DIR}/backups"

# Ensure backup directory exists
mkdir -p "$BACKUP_DIR"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/backup_${DB_NAME}_${TIMESTAMP}.sql.gz"

echo "Starting database backup for $DB_NAME..."

# Check if the docker container is running
if [ "$(docker ps -q -f name=^${CONTAINER_NAME}$)" ]; then
    # Run pg_dump inside the container and stream to host gzipped file
    docker exec -i "$CONTAINER_NAME" pg_dump -U "$DB_USER" -d "$DB_NAME" | gzip > "$BACKUP_FILE"
    echo "Backup completed successfully: $BACKUP_FILE"
else
    echo "Error: Container '$CONTAINER_NAME' is not running." >&2
    exit 1
fi
