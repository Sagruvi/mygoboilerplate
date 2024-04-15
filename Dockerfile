# Используем официальный образ Go как базовый
FROM golang:1.20-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Скачиваем все зависимости
RUN go mod init main && \
    go mod tidy
RUN mkdir proxy
COPY . ./proxy
# Собираем приложение
RUN go build -o main ./proxy/cmd
COPY ./migrations .
COPY config/.env .

# Открываем порт 8080
EXPOSE 8080
EXPOSE 9090

# Запускаем приложение
CMD ["./main"]
