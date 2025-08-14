FROM golang:1.24.4-alpine

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /bot cmd/main.go

RUN apk add --no-cache git ca-certificates tzdata && update-ca-certificates

CMD ["/bot"]