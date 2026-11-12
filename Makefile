SHELL := /bin/bash
.ONESHELL:
.DEFAULT_GOAL := test

.PHONY: test clean run-app stop

test:
	@set -e
	trap 'docker compose down -v' EXIT

	docker compose down -v
	docker compose build test
	docker compose up --wait --wait-timeout 30 db

	docker compose run --rm test sh -c \
		'sqlc generate && templ generate && go build ./... && go test ./... -v'

clean:
	docker compose down -v
	
run-app:
	@set -e
	trap 'docker compose down -v' EXIT

	docker compose down -v
	docker compose build app
	docker compose up --wait --wait-timeout 30 db

	docker compose up app

stop:
	docker compose down
