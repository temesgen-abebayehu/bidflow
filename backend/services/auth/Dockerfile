# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire backend directory
# NOTE: Run this from the 'backend' directory: docker build -f services/auth/Dockerfile .
COPY . .

# Initialize workspace and add modules
# Re-init workspace to ensure clean state
RUN rm -f go.work go.work.sum
RUN go work init
RUN go work use ./common ./services/auth

# Download dependencies
WORKDIR /app/services/auth
RUN go mod download

# Build the application
RUN go build -o auth-service .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/auth/auth-service .

EXPOSE 8080

CMD ["./auth-service"]

# Run from backend/ directory
# docker build -f services/auth/Dockerfile -t auth-service .