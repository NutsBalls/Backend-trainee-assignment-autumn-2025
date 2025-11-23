include .env

run:
	go run ./cmd/pr

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
