POSTGRES_DSN ?= "postgres://shortener:postgres-password@localhost:5432/shortener?sslmode=disable"

migrate-up:
	migrate -path=./scripts/migrations -database ${POSTGRES_DSN} up
.PHONY: migrate-up

migrate-down:
	migrate -path=./scripts/migrations -database ${POSTGRES_DSN} down
.PHONY: migrate-down

build-app:
	go build -o build/app cmd/app/main.go
.PHONY: build-app
