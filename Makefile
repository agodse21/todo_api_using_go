include .env

up: 
	@echo "Starting mongodb containers"
	docker-compose up --build -d --remove-orphans

down:	
	@echo "Stopping mongodb containers"
	docker-compose down


build:
	go build -o $(PROJECT_NAME) ./cmd/api/


start:
	./$(PROJECT_NAME)


restart: build start
   

