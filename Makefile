.PHONY: build run shutdown postgres create-migrate migrate

build:
	docker build --tag dnevsky/tg-bot-service .

run:
	docker-compose up -d tg-bot-service

shutdown:
	docker-compose down

postgres:
	docker-compose up -d postgres

create-migrate:
	migrate create -ext sql -dir ./repos/postgres/migrations -seq <name>

migrate:
	docker-compose run migrate-postgres