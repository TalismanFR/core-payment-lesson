# docker network for containers
net = app-net-1

# creates docker network
network:
	bash scripts/docker_network.sh $(net)

# starts vault container and initialize it with values
# after execution copy VAULT_TOKEN to .env
# if you want to test grpc server, copy one of the uuid and paste as terminalId in test/grpc_client/main.go
vault: network
	cd deploy_vault; \
	bash vault.sh $(net) && \
	go run main.go --file=terminals.json

# builds pay-service
build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

# runs pay-service and postgres in containers
run: build network
	docker-compose up --remove-orphans --build pay

# stops pay-service, postgres and vault
stop:
	docker compose down
	docker stop vault

# clears builds, network, pulled images
clear:
	sudo rm -rf .bin
	docker network rm $(net)
	docker image rm postgres:14.3-alpine3.16
	docker image rm vault:1.10.4
	docker image rm pay-service

test.integration:
	go test -tags=integration -v ./...