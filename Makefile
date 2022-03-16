start:
	docker-compose up -d

start-build:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	make stop && make start


restart-build:
	make stop && make start-build

dev:
	APP_HOST=:8888 fresh -c other_runner.conf

install:
	export GOPRIVATE=gitlab.finema.co/finema/* && git config --global url."git@gitlab.finema.co:".insteadOf "https://gitlab.finema.co/" && go get

logs:
	 docker logs -f api

migrate:
	docker-compose up -d --build migration

test:
	echo 'test ./...'

download-modules:
	echo 'go mod download'

test-e2e:
	go test --tags=e2e ./...
