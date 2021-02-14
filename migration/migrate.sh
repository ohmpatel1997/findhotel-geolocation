#!/bin/bash

if [ $2 = "down" ]; then
    echo "Cannot run down, use down-to instead"
    exit 1
fi

# Build golang package
go build -o goose *.go

# run commands
./goose -dir "${1}" postgres "dbname=${DB_NAME} user=${DB_USER} password=${DB_PASSWD} host=${HOST} sslmode=disable" $2 $3
