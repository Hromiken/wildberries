.PHONY: build up down logs ps restart curl-order

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker logs -f app

ps:
	docker ps

restart: down up

curl-order:
	curl -X POST http://localhost:8080/new_order \
	  -H "Content-Type: application/json" \
	  -d @example_order.json
