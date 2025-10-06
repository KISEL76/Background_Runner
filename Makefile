APP_NAME=queue-svc
DOCKER_IMAGE=queue-svc:latest
INTERNAL_DIR=./internal/

all: build run

build:
	docker build -t $(DOCKER_IMAGE) .

run: 
	docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE) || true

test:
	go test -v $(INTERNAL_DIR)store $(INTERNAL_DIR)queue $(INTERNAL_DIR)worker $(INTERNAL_DIR)integration_test

test-race: 
	go test -race ./...

vet:
	go vet ./...

clean:
	docker rmi $(DOCKER_IMAGE) || true

.PHONY: build run clean test test-race all
