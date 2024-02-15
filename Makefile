build:
	docker-compose build weather-bot-app

run:
	docker-compose up weather-bot-app

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up

