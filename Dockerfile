FROM golang:1.24-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o chat_service ./cmd/main.go 


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/chat_service  .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./chat_service", "--config=./config/local.yaml"]