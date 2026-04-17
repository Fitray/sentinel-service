include .env
export
export PROJECT_ROOT=${shell pwd}
export LOGGER_PATH=${PROJECT_ROOT}/out/logger

run:
	@go run ${PROJECT_ROOT}/cmd/app/main.go

run-python:
	@${PROJECT_ROOT}/.venv/bin/python ${PROJECT_ROOT}/internal/python/main.py ${city}