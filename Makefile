###############################
# Local
###############################
dot-env:
	cp .env.template .env
	cp .env.template .env.docker

build-migrate:
	go build -o bin/migrate cmd/migrate/main.go

build-server:
	go build -o bin/server cmd/server/main.go

run-migrations: build-migrate
	./bin/migrate

run-server: build-server
	./bin/server

run-tests:
	go test -v ./...

###############################
# Docker
###############################
docker-build:
	docker-compose build --no-cache

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down