# Stage 1: Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache bash git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/avito main.go
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/avito .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/migrate.sh .
RUN apk add --no-cache bash postgresql-client
RUN chmod +x migrate.sh
ENV DATABASE_URL=postgres://user:pass@db:5432/yourdb?sslmode=disable
CMD ["sh", "-c", "/app/migrate.sh && /app/avito"]
