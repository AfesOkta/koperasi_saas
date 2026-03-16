.PHONY: dev build run test migrate seed docker-up docker-down lint clean

# ─── Development ───
dev:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

# ─── Testing ───
test:
	go test ./... -v -cover

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# ─── Database ───
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

seed:
	go run cmd/seeder/main.go

# ─── Docker ───
docker-up:
	docker compose -f ../infra/docker/docker-compose.yml up -d

docker-down:
	docker compose -f ../infra/docker/docker-compose.yml down

# ─── Code Quality ───
lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

# ─── Dependencies ───
deps:
	go mod tidy
	go mod download

# ─── Clean ───
clean:
	rm -rf bin/ coverage.out coverage.html
