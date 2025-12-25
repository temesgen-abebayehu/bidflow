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
# -ldflags="-w -s" reduces binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o notification-service .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/services/notification/notification-service .

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

EXPOSE 8084

CMD ["./notification-service"]

