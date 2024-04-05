prepare:
	cp .env.example ./.deploy/.env

start:
	docker compose -f ./.deploy/docker-compose.yml up

prepare-raw:
	cp .env.example ./.env && yarn install --frozen-lockfile && yarn build

start-raw:
	yarn start
