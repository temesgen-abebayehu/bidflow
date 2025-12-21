# Generate gRPC code from .proto files
proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

# Spin up infrastructure (Postgres, Kafka)
up:
	docker-compose up -d

# Stop infrastructure
down:
	docker-compose down

.PHONY: proto up down