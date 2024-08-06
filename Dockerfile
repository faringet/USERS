FROM golang:1.21.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY config/conf.yaml config/conf.yaml

ENV ENV=docker

EXPOSE 3003

CMD ["./main"]
