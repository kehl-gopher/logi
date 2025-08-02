DOCKER_COMPOSE = ./infra/docker-compose.yml
.DEFAULT_GOAL := docker-up
POSTGRESQL_URL='postgres://logi:logi123@localhost:5432/logi?sslmode=disable'
MIGRATE_PATH = migrations/

docker-up:
	@docker compose -f $(DOCKER_COMPOSE) --env-file ./.env up

docker-down:
	@docker compose -f $(DOCKER_COMPOSE) down

.PHONY: run compile
run: compile
	./main

compile:
	@go build -o main

.PHONY: test
test:
	@go test ./tests/...
	
.PHONY: fmt add
fmt:
	@gofmt -l .
	@gofmt -w .

add: fmt
	git add .

.PHONY: clean
clean:
	@rm ./main

.PHONY: migrate-create migrate-up migrate-down
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $$name

migrate-up: 
	migrate -database $(POSTGRESQL_URL) -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database $(POSTGRESQL_URL) -path $(MIGRATE_PATH) down

.PHONY: drop
drop:
	./main -drop