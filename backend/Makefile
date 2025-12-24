# Generate gRPC code from .proto files
proto:
	protoc --proto_path=proto --go_out=proto/pb --go_opt=paths=source_relative --go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative proto/*.proto

# Spin up infrastructure (Postgres, Kafka)
up:
	docker-compose up -d

# Stop infrastructure
down:
	docker-compose down

# Run go mod tidy for all modules
tidy:
	cd common && go mod tidy
	cd proto && go mod tidy
	cd services/auth && go mod tidy
	cd services/auction && go mod tidy
	cd services/bidding && go mod tidy
	cd services/notification && go mod tidy

# Build all services
build:
	cd services/auth && go build -o ../../bin/auth main.go
	cd services/auction && go build -o ../../bin/auction main.go
	cd services/bidding && go build -o ../../bin/bidding main.go
	cd services/notification && go build -o ../../bin/notification main.go

# Run all tests
test:
	go test ./common/...
	go test ./services/auth/...
	go test ./services/auction/...
	go test ./services/bidding/...
	go test ./services/notification/...

# Clean build artifacts
clean:
	rm -rf bin/

.PHONY: proto up down tidy build test clean