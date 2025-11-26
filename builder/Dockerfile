FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker cmd/worker/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install dependencies for building
RUN apk --no-cache add \
    git \
    docker-cli \
    nodejs \
    npm \
    go \
    php \
    composer

# Copy binary
COPY --from=builder /app/bin/worker .

CMD ["./worker"]

