build:
	docker compose build --no-cache lay-deputy-hub

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --no-cache --parallel lay-deputy-hub json-server

clean:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml down --remove-orphans

dev-up: clean build-dev npm
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up lay-deputy-hub npm json-server

dev-up-sirius: clean build-dev npm
	docker compose -f docker-compose.yml -f docker/docker-compose.sirius.yml up lay-deputy-hub npm

down:
	docker compose down --remove-orphans

go-lint:
	docker compose run --rm go-lint

npm:
	docker compose run --rm npm

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs

unit-test: test-results
	docker compose run --rm test-runner

up: clean build
	docker compose up -d --wait lay-deputy-hub