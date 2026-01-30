init: up \
	  build \
	  run-migrations

up:
	docker compose up -d

build:
	docker compose build

down:
	docker compose down

# команда создания миграции name=create-chat-table
create-migration:
	  docker compose run --rm migration sh -c 'goose -dir migrations create $(name) sql'

run-migrations:
	  docker compose run --rm migration sh -c 'goose -dir migrations postgres "$$DB_URL" up'