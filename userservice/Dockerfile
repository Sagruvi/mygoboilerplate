FROM golang:1.22.0-alpine3.19
LABEL authors="sagruvi"


WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o userservice ./cmd/
CMD ["./userservice"]
EXPOSE 15002
