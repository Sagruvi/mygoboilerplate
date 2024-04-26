FROM golang:1.22.0-alpine3.19
LABEL authors="sagruvi"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o proxy ./cmd/

CMD [ "./proxy" ]

EXPOSE 8080