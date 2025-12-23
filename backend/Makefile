# Generate gRPC code from .proto files
proto:
	protoc --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative auction.proto bidding.proto

# Spin up infrastructure (Postgres, Kafka)
up:
	docker-compose up -d

# Stop infrastructure
down:
	docker-compose down

.PHONY: proto up down