# Используем официальный образ Go как базовый
FROM golang:1.20-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Устанавливаем gcc, unzip и другие зависимости для сборки
RUN apk add --no-cache gcc musl-dev

# Копируем go.mod и go.sum (если они есть) для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Установка swagger
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.0

# Копируем исходники приложения в рабочую директорию
COPY . .

# Проверка наличия swagger.yaml
RUN ls -la ./swagger/

# Собираем приложение
RUN go build -o  main .

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main
# Копируем swagger в корень
COPY --from=builder /app/swagger /swagger

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
