# Этап сборки
FROM golang:latest AS builder

#RUN apk update && apk add --no-cache git gcc musl-dev

# Копируем исходный код
WORKDIR /app
COPY . .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-extldflags "-static"' -o music-service ./cmd/server/

# Этап рантайма
FROM alpine:3.18

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Рабочая директория
WORKDIR /app

# Копируем бинарник и конфиг
COPY --from=builder --chown=appuser:appgroup /app/music-service .
COPY --from=builder --chown=appuser:appgroup /app/config.yaml .

# Переключаемся на непривилегированного пользователя
USER appuser

# Проверяем что файл существует и исполняемый
RUN ls -la && [ -f music-service ] && [ -x music-service ]

# Запускаем приложение
CMD ["./music-service"]