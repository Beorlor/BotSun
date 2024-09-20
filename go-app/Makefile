# Makefile

.PHONY: up down build logs

up:
	docker-compose up -d go-app timescaledb

down:
	docker-compose down

build:
	docker-compose build go-app timescaledb

logs:
	docker-compose logs -f go-app timescaledb
