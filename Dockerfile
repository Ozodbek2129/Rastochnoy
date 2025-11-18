FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# MAIN FILE: /app/cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]
