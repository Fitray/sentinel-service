include .env
export
export PROJECT_ROOT=${shell pwd}
export LOGGER_PATH=${PROJECT_ROOT}/out/logger
export POSTGRES_HOST=localhost

run:
	@go run ${PROJECT_ROOT}/cmd/app/main.go

run-python:
	@${PROJECT_ROOT}/.venv/bin/python ${PROJECT_ROOT}/internal/python/main.py ${city}

up-postgres:
	@mkdir -p ${PROJECT_ROOT}/out/pgdata
	docker compose up -d sentinel-postgres

down-postgres:
	@docker compose down sentinel-postgres

cleanup-postgres:
	@make down-postgres && \
	rm -rf out/pgdata

migrate-create:
	@migrate create -ext sql -dir ${PROJECT_ROOT}/migrations -seq $(seq)

migrate-action:
	docker compose run --rm postgres-migrate  \
	-path migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@sentinel-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	$(action)

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

env-port-forward:
	@sudo systemctl stop postgresql
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

connection-line:
	@echo postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable