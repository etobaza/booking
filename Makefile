.PHONY: test up down run-tests

up:
	docker-compose up -d

down:
	docker-compose down

run-tests:
	docker-compose run --rm app go test ./...

test: up run-tests down
