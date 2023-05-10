FROM golang:1.20-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o tg-bot-service cmd/main.go
 
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/tg-bot-service /app/tg-bot-service
COPY --from=builder /app/data /app/data

WORKDIR /app

CMD ["./tg-bot-service"]