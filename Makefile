.PHONY: make install run down clean test

make:
	@echo "Available commands:"
	@echo "  make install   Install dependencies and tools"
	@echo "  make run       Run API + Redis in Docker"
	@echo "  make down      Stop containers"
	@echo "  make clean     Remove containers + volumes"
	@echo "  make test      Run unit tests"

install:
	@which docker >/dev/null 2>&1 || { echo "Docker not found. Install Docker: https://docs.docker.com/get-docker/"; exit 1; }
	@which docker-compose >/dev/null 2>&1 || { echo "docker-compose not found. Install: https://docs.docker.com/compose/install/"; exit 1; }
	@which go >/dev/null 2>&1 || { echo "Go not found. Install: https://go.dev/dl/"; exit 1; }
	go mod tidy

run:
	docker-compose up --build

down:
	docker-compose down

clean:
	docker-compose down -v

test:
	go test ./...
