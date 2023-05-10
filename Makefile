.PHONY: build run shutdown

build:
	docker build --tag dnevsky/tg-bot-service .

run:
	docker-compose up -d tg-bot-service

shutdown:
	docker-compose down