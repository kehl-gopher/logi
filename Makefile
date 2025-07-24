
DOCKER_COMPOSE = ./infra/docker-compose.yml

docker-up:
	docker compose -f $(DOCKER_COMPOSE) --env-file ./.env up

docker-down:
	docker compose -f $(DOCKER_COMPOSE) down

run: compile
	./main

compile:
	go build -o main

test:
	go test ./tests/...
	
fmt:
	gofmt -l .
	gofmt -w .

clean:
	rm ./main