init:
	go mod tidy && go mod download

run:
	make init && go run src/main.go

build:
	make init && go build -ldflags="-s -w" -o ./dist/app src/main.go

dev:
	make init && mkdir -p dist/tmp && air server



start:
	cp -n .env.example ./.deploy/.env && docker compose -f ./.deploy/docker-compose.yml up

start-raw:
	cp -n .env.example ./.env && make run
