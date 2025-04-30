FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Устанавливаем swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
COPY . .

# Генерируем документацию и скачиваем зависимости
RUN cd /app && \
    swag init -g cmd/main.go -o internal/docs && \
    go mod download && \
    go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/internal/docs ./internal/docs

EXPOSE 8080
CMD ["./main"]