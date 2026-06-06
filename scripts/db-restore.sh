#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"

# Load environment variables from .env.production if it exists
if [ -f "$PARENT_DIR/.env.production" ]; then
    export $(grep -v '^#' "$PARENT_DIR/.env.production" | xargs)
fi

# Configuration with defaults
DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-chatsphere}
CONTAINER_NAME=${CONTAINER_NAME:-chatsphere-postgres-prod}

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <path_to_backup_file.sql.gz>"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Check if the docker container is running
if [ -z "$(docker ps -q -f name=^${CONTAINER_NAME}$)" ]; then
    echo "Error: Container '$CONTAINER_NAME' is not running."
    exit 1
fi

# Confirm action
read -p "WARNING: This will overwrite the database '$DB_NAME'. Are you sure you want to proceed? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Restore cancelled by user."
    exit 0
fi

echo "Restoring database from $BACKUP_FILE..."

# Terminate active connections, drop and recreate DB to guarantee clean slate
echo "Recreating database '$DB_NAME'..."
docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d postgres -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '${DB_NAME}' AND pid <> pg_backend_pid();" || true
docker exec -i "$CONTAINER_NAME" dropdb -U "$DB_USER" --if-exists "$DB_NAME"
docker exec -i "$CONTAINER_NAME" createdb -U "$DB_USER" "$DB_NAME"

# Restore database
gunzip -c "$BACKUP_FILE" | docker exec -i "$CONTAINER_NAME" psql -U "$DB_USER" -d "$DB_NAME"

echo "Database restore completed successfully!"
