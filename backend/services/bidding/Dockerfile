# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire backend directory
# NOTE: Run this from the 'backend' directory: docker build -f services/bidding/Dockerfile .
COPY . .

# Initialize workspace and add modules
# Re-init workspace to ensure clean state inside the container
RUN rm -f go.work go.work.sum
RUN go work init
RUN go work use ./common ./proto ./services/bidding

# Download dependencies
WORKDIR /app/services/bidding
RUN go mod download

# Build the application
RUN go build -o bidding-service .

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/services/bidding/bidding-service .

# Expose HTTP and gRPC ports
EXPOSE 8080
EXPOSE 50051

CMD ["./bidding-service"]
