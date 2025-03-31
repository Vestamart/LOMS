FROM golang:1.23.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o loms-service ./cmd/server

FROM alpine:latest
WORKDIR /app

# Устанавливаем необходимые зависимости для Goose и PostgreSQL клиента
RUN apk add --no-cache postgresql-client

# Копируем бинарник Goose из сборочного образа
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Копируем остальные файлы
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/loms-service .
COPY --from=builder /app/stock-data.json .

EXPOSE 50051

# Запускаем миграции и затем сервис
CMD ["sh", "-c", "goose -dir ./migrations postgres 'postgres://root:root@postgres:5432/loms_db?sslmode=disable' up && ./loms-service"]