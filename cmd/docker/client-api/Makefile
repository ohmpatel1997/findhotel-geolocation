
all: build docker-build run-local

build: docker-build
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o main .

docker-build:
	docker build \
		.

run-local:
	docker-compose build
	docker-compose up
