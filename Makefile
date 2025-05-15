.PHONY: build run test migrate

build:
	go build -o bin/loader ./cmd/loader

run:
	go run ./cmd/loader/main.go

test:
	go test -v ./...

migrate-up:
	goose -dir migrations postgres "user=postgres dbname=mydb sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "user=postgres dbname=mydb sslmode=disable" down

docker-build:
	docker build -t data-loader .

docker-run:
	docker run -v ./configs:/app/configs data-loader
