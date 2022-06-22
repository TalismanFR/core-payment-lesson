docker.network:
	docker network create app-net-1

run.app: docker.network
	docker compose up

run.vault: docker.network
	make ./other/vault/Makefile net=app-net-1

test.integration:
	go test  -tags=integration -v ./...

build:
	go build -o ./.bin/app ./cmd/main.go

run: build
	docker-compose up --remove-orphans app

