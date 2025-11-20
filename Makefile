include .env

run:
	go run ./cmd/pr

migrate-create:
	goose -dir migrations -s create $(name) sql

migrate-up:
	goose -dir migrations postgres "$(DB_CONN)" up

migrate-down:
	goose -dir migrations postgres "$(DB_CONN)" down

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
