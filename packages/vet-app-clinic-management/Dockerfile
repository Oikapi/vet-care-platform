FROM golang:1.24.2-bullseye AS builder

WORKDIR /app

# 1. Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 2. Установка swag
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:$PATH"

# 3. Копируем остальной код
COPY . .

# 4. Генерируем swagger
RUN swag init -g cmd/main.go -o internal/docs

# 5. Сборка приложения
RUN go build -o app ./cmd/main.go

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]

