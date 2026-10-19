SHELL := /bin/bash
.ONESHELL:
.DEFAULT_GOAL := test

.PHONY: test clean run stop test-api

test:
	@set -e
	trap 'docker compose down -v' EXIT

	docker compose down -v
	docker compose build test
	docker compose up --wait --wait-timeout 30 db

	docker compose run --rm test sh -c \
		'sqlc generate && go build ./... && go test ./... -v'

clean:
	docker compose down -v
	
run:
	docker compose up --build app

stop:
	docker compose down
	

test-api:
	@set -e
	trap 'docker compose -p alfajorboxd-test-api down -v' EXIT

	docker compose -p alfajorboxd-test-api down -v
	docker compose -p alfajorboxd-test-api up -d --build app

	echo "Esperando que la API esté disponible..."

	for i in {1..60}; do
		if curl -fsS http://localhost:8080/alfajores >/dev/null 2>&1; then
			break
		fi

		if [ "$$i" -eq 60 ]; then
			echo "ERROR: La API no respondió"
			docker compose -p alfajorboxd-test-api logs --tail=80 app db
			exit 1
		fi

		sleep 2
	done

	bash request.sh

