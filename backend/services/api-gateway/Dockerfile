FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.work go.work.sum ./
COPY services/api-gateway/go.mod ./services/api-gateway/
COPY common/go.mod common/go.sum ./common/

# Adjust paths if we are copying from root context
# Since the context in docker-compose is likely '.', we need to copy everything relevant.
# Assuming standard go build flow within the monorepo context or service context.
# Looking at auth-service Dockerfile in docker-compose: dockerfile: services/auth/Dockerfile, context: .
# So we have access to the whole repo.

COPY services/api-gateway ./services/api-gateway
COPY common ./common

WORKDIR /app/services/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/services/api-gateway/api-gateway .
# COPY --from=builder /app/.env . # Optional, usually passed by docker-compose

EXPOSE 8080

CMD ["./api-gateway"]
