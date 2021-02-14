export
HOST=localhost
DB_SSL_MODE=disable
DB_USER=ohmpatel
DB_PASSWD=ohmpatel
DB_NAME=postgres

migrate-status: build
	$$(sh migrate.sh geolocation status)
migrate-down: build
	$$(sh migrate.sh geolocation down)
migrate-up: build
	$$(sh migrate.sh geolocation up)

build:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -o main .