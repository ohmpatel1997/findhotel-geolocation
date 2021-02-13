#!/bin/bash

if [ $2 = "down" ]; then
    echo "Cannot run down, use down-to instead"
    exit 1
fi

# Build golang package
go build -o goose *.go

# run commands
./goose -dir "${1}" postgres "dbname=${RDS_DB_NAME} user=${RDS_USERNAME} password=${RDS_PASSWORD} host=${RDS_HOSTNAME} sslmode=disable" $2 $3
