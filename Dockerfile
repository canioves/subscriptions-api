FROM golang:1.25-alpine AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/main.go

COPY migrations ./migrations

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]