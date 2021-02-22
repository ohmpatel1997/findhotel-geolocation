export
HOST=localhost
PORT=5432
DB_SSL_MODE=disable
DB_USER=ohmpatel
DB_PASSWD=ohmpatel
DB_NAME=postgres
FILE_NAME=data_dump.csv

all: build run

build:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o main .

run: build
	go run main.go

