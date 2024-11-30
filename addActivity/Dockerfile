# Шаг 1: Используем базовый образ Go на основе Alpine для сборки приложения
FROM golang:alpine AS builder
LABEL stage=gobuilder

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /build

# Копируем зависимости и скачиваем их в отдельный слой (оптимизация кэша)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .
RUN go build -ldflags="-s -w" -o /app/main .

# Шаг 2: Используем минимальный Alpine образ для запуска приложения
FROM alpine
RUN apk add --no-cache ca-certificates postgresql-client

WORKDIR /app
# Копируем собранный бинарный файл из builder-образа в конечный образ
COPY --from=builder /app/main .
# Копируем скрипт ожидания
COPY wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

EXPOSE 8080
CMD ["/app/wait-for-postgres.sh", "db", "/app/main"]
