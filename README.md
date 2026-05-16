# Customer API

`go-customer-api` is a Go REST API for basic customer management.  
It uses Gorilla Mux for routing, PostgreSQL for persistence, and Gormigrate to run schema/data migrations automatically on startup.

## Features

- CRUD endpoints for customers
- Request validation for customer IDs and payload fields
- Auto-migration and seed data at application startup
- Swagger/OpenAPI documentation

## Tech Stack

- Go
- Gorilla Mux
- PostgreSQL
- GORM + Gormigrate
- Swaggo (`swagger` UI)

## Project Structure

```text
controllers/   HTTP handlers
db/            DB connection and runtime migrations
docs/          Generated Swagger files
models/        Data models
routes/        Route registration and middleware
script/        Local DB bootstrap scripts
main.go        Application entrypoint
```

## Prerequisites

- Go 1.24+
- PostgreSQL (local install or Docker)

## Environment Variables

Copy `.env.example` to `.env` and fill values:

```env
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=yourdbname
DB_HOST=localhost
DB_PORT=5432
```

The app requires all five variables.

## Run PostgreSQL With Docker

Use either option:

1. `docker compose up -d postgres`
2. `./script/create_postgres_db.sh`

The script can read values from `.env` and creates the `customers` table if missing.

## Run the API

```sh
go mod tidy
go run main.go
```

Server URL: `http://localhost:8000`

On startup, the app:

- Connects to PostgreSQL
- Runs Gormigrate migrations in `db/db.go`
- Seeds sample customers if they do not already exist

## API Endpoints

- `GET /api/customers` - list all customers
- `GET /api/customers/{id}` - get one customer
- `POST /api/customers` - create customer
- `PUT /api/customers/{id}` - update customer
- `DELETE /api/customers/{id}` - delete customer

Customer payload:

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com"
}
```

## Swagger Docs

After server start:

`http://localhost:8000/swagger/index.html`

Regenerate docs after changing handler annotations:

```sh
swag init
```

Install `swag` if needed:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
