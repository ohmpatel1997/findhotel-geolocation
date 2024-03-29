#!/bin/bash


# Build golang package
go build -o goose *.go

# run commands
./goose -dir "${1}" postgres "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWD} host=${HOST} sslmode=${DB_SSL_MODE}" $2 $3
