run:
	@echo "Start app"
	go run ./cmd/gophermart/main.go

build:
	@echo "Build app"
	go build -o ./cmd/gophermart/main.go
all:
	@echo "Start build app"

.PHONY: cover
cover:
	go test -short -count=1 -coverprofile=coverage.out ./internal/... ./cmd/...
	go tool cover -html=coverage.out
	rm coverage.out
migrate:
	@echo "Start migrate up"
	migrate -path migrations -database "postgres://user:password@localhost:5432/loyalty?sslmode=disable" up

down:
	@echo "Start migrate down"
	migrate -path migrations -database "postgres://user:password@localhost:5432/loyalty?sslmode=disable" down

#// migrate create -ext sql -dir migrations -seq create_users_table
#  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# migrate -path migrations -database "postgres://user:password@localhost:5432/loyalty?sslmode=disable" up

# mockery --name=Repository --output=mocks --outpkg=mocks