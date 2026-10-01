build:
	docker compose build --no-cache lay-deputy-hub

build-dev:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml build --no-cache --parallel lay-deputy-hub json-server

clean:
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml down --remove-orphans

compile-assets:
	docker compose run --rm npm build

dev-up: clean build-dev
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up lay-deputy-hub watch-assets json-server

up: clean compile-assets build
	docker compose up -d --wait lay-deputy-hub

dev-up-sirius: clean build-dev npm
	docker compose -f docker-compose.yml -f docker/docker-compose.sirius.yml up lay-deputy-hub npm

down:
	docker compose down --remove-orphans

go-lint:
	docker compose run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs

unit-test: test-results
	docker compose run --rm test-runner