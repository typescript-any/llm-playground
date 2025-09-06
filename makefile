.PHONY: dev migrate-up migrate-down db-up db-down 

# Start hot-reload dev server using Air
dev:
	@air 

# Run DB migration up
migrate-up: 
	go run cmd/migrate/main.go -direction up

# Run DB migration down
migrate-down:
	go run cmd/migrate/main.go -direction down

# Up database using docker-compose
db-up:
	docker-compose up -d

db-down:
	docker-compose down