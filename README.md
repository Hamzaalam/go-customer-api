# Customer API

## Run the API

```sh
go run main.go
```

The API runs on `http://localhost:8000`.

## Create PostgreSQL Database With Docker

```sh
./script/create_postgres_db.sh
```

The script reads `DB_USER`, `DB_PASSWORD`, `DB_NAME`, and `DB_PORT` from `.env` when present. If they are missing, it uses local development defaults and prints the values to add to `.env`.

## Swagger Documentation

Swagger UI is available after the server starts:

```text
http://localhost:8000/swagger/index.html
```

Regenerate the Swagger docs after changing route annotations:

```sh
swag init
```

If `swag` is not installed:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
