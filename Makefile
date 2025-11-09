include .env
export

# Makefile raíz de gestor-one

# Variables
DB_DSN ?= $(DB_DSN)
MIGRATIONS_DIR = migrations

# Comandos

## Correr migraciones hacia adelante
migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

## Revertir migraciones (un paso)
migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1

## Crear una nueva migración
# Uso: make new-migration name=create_table_users
new-migration:
	@if [ -z "$(name)" ]; then \
		echo "Falta el parámetro 'name'. Ej: make new-migration name=init"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## Ver el estado actual de migraciones
migrate-version:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version

## Build local
build:
	@go build -o gestor-one ./cmd/server

## Correr localmente
run:
	@go run ./cmd/server

## Limpieza
clean:
	@rm -f gestor-one
