# Dockerfile
FROM golang:1.23.2-alpine AS builder

# Установка зависимостей
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода и сборка сервера
COPY . .
RUN go build -o Golang-Proxy-Youtube ./src/server/main.go

# Финальный минимальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/Golang-Proxy-Youtube /app/Golang-Proxy-Youtube

# Открытие порта для gRPC
EXPOSE 50051

# Команда для запуска сервера
CMD ["/app/Golang-Proxy-Youtube"]
