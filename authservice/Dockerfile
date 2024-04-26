FROM golang:1.22.0-alpine3.19

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod tidy
COPY . .

RUN go build -o auth ./cmd/

CMD ["./auth"]
EXPOSE 8080
