.PHONY: run build migrate-up migrate-down seed swagger tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

create-migration:
	go run ./cmd/migration create migration_name

migrate-up:
	go run ./cmd/migration up

migrate-down:
	go run ./cmd/migration down

seed:
	go run ./cmd/seeder

swagger:
	swag init -g cmd/api/main.go --output docs


