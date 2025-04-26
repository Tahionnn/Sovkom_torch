# Используем официальный образ Go для сборки
FROM golang:1.23.8 AS builder

WORKDIR /app

# Копируем файлы проекта и зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Финальный минимальный образ
FROM debian:latest

WORKDIR /app
COPY --from=builder /app/main .

# Устанавливаем CA-сертификаты
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Даем права на выполнение
RUN chmod +x main

# Запуск
CMD ["./main"]
