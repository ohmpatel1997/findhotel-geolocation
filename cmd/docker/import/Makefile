export
HOST=ec2-52-7-168-69.compute-1.amazonaws.com
PORT=5432
DB_SSL_MODE=require
DB_USER=hzlujbndpdgarb
DB_PASSWD=7f41485aa7886d8f3c541407c504ebf17769f4dc88d6300fc32829e4103b5091
DB_NAME=d5e4jb9muq5qn2
FILE_NAME=data_dump.csv

all: build run

build:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o main .

run: build
	go run main.go

