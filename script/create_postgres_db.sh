#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ENV_FILE:-$ROOT_DIR/.env}"

if [[ -f "$ENV_FILE" ]]; then
	set -a
	source "$ENV_FILE"
	set +a
fi

DB_USER="${DB_USER:-customer_user}"
DB_PASSWORD="${DB_PASSWORD:-customer_password}"
DB_NAME="${DB_NAME:-customer_api}"
DB_PORT="${DB_PORT:-5432}"

POSTGRES_CONTAINER_NAME="${POSTGRES_CONTAINER_NAME:-customer-api-postgres}"
POSTGRES_IMAGE="${POSTGRES_IMAGE:-postgres:16-alpine}"
POSTGRES_VOLUME="${POSTGRES_VOLUME:-${POSTGRES_CONTAINER_NAME}-data}"
POSTGRES_ADMIN_DB="${POSTGRES_ADMIN_DB:-postgres}"

if ! command -v docker >/dev/null 2>&1; then
	echo "docker is required but was not found in PATH" >&2
	exit 1
fi

if ! docker container inspect "$POSTGRES_CONTAINER_NAME" >/dev/null 2>&1; then
	echo "Creating PostgreSQL container: $POSTGRES_CONTAINER_NAME"
	docker run -d \
		--name "$POSTGRES_CONTAINER_NAME" \
		-e POSTGRES_USER="$DB_USER" \
		-e POSTGRES_PASSWORD="$DB_PASSWORD" \
		-e POSTGRES_DB="$POSTGRES_ADMIN_DB" \
		-p "$DB_PORT:5432" \
		-v "$POSTGRES_VOLUME:/var/lib/postgresql/data" \
		"$POSTGRES_IMAGE" >/dev/null
elif [[ "$(docker inspect -f '{{.State.Running}}' "$POSTGRES_CONTAINER_NAME")" != "true" ]]; then
	echo "Starting PostgreSQL container: $POSTGRES_CONTAINER_NAME"
	docker start "$POSTGRES_CONTAINER_NAME" >/dev/null
else
	echo "PostgreSQL container is already running: $POSTGRES_CONTAINER_NAME"
fi

echo "Waiting for PostgreSQL to accept connections..."
until docker exec "$POSTGRES_CONTAINER_NAME" pg_isready -U "$DB_USER" -d "$POSTGRES_ADMIN_DB" >/dev/null 2>&1; do
	sleep 1
done

echo "Creating database if missing: $DB_NAME"
docker exec \
	-e PGPASSWORD="$DB_PASSWORD" \
	"$POSTGRES_CONTAINER_NAME" \
	psql -U "$DB_USER" -d "$POSTGRES_ADMIN_DB" -v ON_ERROR_STOP=1 -tc \
	"SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" |
	grep -q 1 ||
	docker exec \
		-e PGPASSWORD="$DB_PASSWORD" \
		"$POSTGRES_CONTAINER_NAME" \
		psql -U "$DB_USER" -d "$POSTGRES_ADMIN_DB" -v ON_ERROR_STOP=1 -c \
		"CREATE DATABASE \"$DB_NAME\""

echo "Creating customers table if missing"
docker exec \
	-e PGPASSWORD="$DB_PASSWORD" \
	"$POSTGRES_CONTAINER_NAME" \
	psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -c \
	'CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE
	);'

cat <<EOF
PostgreSQL is ready.

Container: $POSTGRES_CONTAINER_NAME
Database:  $DB_NAME
User:      $DB_USER
Port:      $DB_PORT

Use these values in .env:
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
DB_HOST=localhost
DB_PORT=$DB_PORT
EOF
