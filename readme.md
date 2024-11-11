# Shortener API

The Shortener API is a service that makes it easy to create short, shareable URLs and manage redirections. In this guide, you’ll find everything you need to set up the service, run migrations, and explore the API with Swagger documentation.

## What You’ll Need

Before you get started, make sure you’ve installed these tools:

- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) - for generating OpenAPI client/server code.
- [golang-migrate](https://github.com/golang-migrate/migrate) - for managing database migrations.

## Getting Started

### Migrations

Database migrations set up and manage the schema for the Shortener API.

#### Running Migrations

To apply migrations and set up the database schema:

```bash
# default

POSTGRES_DSN="postgres://shortener:postgres-password@localhost:5432/shortener?sslmode=disable"

# you can change it and use like this

POSTGRES_DSN="postgres..." make migrate-up
```

```bash
make migrate-up
```

#### Dropping Migrations
To roll back the migrations:

```bash
make migrate-down
```

### Docker
Docker is used to containerize the Shortener API, making it easy to set up and run.

#### Starting Docker Compose
To build and run the containers in detached mode:

```bash
docker compose up --build -d
```

### Stopping Docker Compose
To stop and remove the containers:

```bash
docker compose down
```

### API Documentation
Swagger documentation is available to interact with the API and view available endpoints.

To access the Swagger docs, open your browser and go to:
```bash
http://localhost:8000/docs
```