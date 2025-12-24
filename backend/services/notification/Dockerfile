# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire backend directory
# NOTE: Run this from the 'backend' directory: docker build -f services/notification/Dockerfile .
COPY . .

# Initialize workspace and add modules
# Re-init workspace to ensure clean state
RUN rm -f go.work go.work.sum
RUN go work init
RUN go work use ./common ./services/notification

# Download dependencies
WORKDIR /app/services/notification
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service .

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/services/notification/notification-service .
# Copy config if needed, or rely on env vars
# COPY --from=builder /app/services/notification/config ./config

EXPOSE 8084

CMD ["./notification-service"]
