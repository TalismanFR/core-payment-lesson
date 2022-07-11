# docker network for containers
net = app-net-1

# creates docker network
network:
	bash scripts/docker_network.sh $(net)

# starts vault container and initialize it with values
# after execution copy VAULT_TOKEN to .env
# if you want to test grpc server, copy one of the uuid and paste as terminalId in test/grpc_client/main.go
start.vault: network
	cd deploy_vault; \
	bash vault.sh $(net) && \
	go run main.go --file=terminals.json

# stop vault container
stop.vault:
	docker stop vault

# builds pay-service
build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

# runs pay-service and postgres in containers
run: build network
	docker compose up --remove-orphans --build pay

# stops pay-service, postgres and jaeger
stop:
	docker compose down

# clears builds, network, pulled images
clean:
	sudo rm -rf .bin
	docker network rm $(net)
	docker image rm postgres:14.3-alpine3.16
	docker image rm vault:1.10.4
	docker image rm pay-service
	docker image rm jaegertracing/all-in-one:1.35

test.integration:
	go test -tags=integration -v ./...