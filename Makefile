run:
	@echo "Запуск приложения..."
	go run ./cmd/gophermart/main.go

build:
	@echo "Сборка приложения..."
	go build -o ./cmd/gophermart/main.go
all:
	@echo "Start build app"
