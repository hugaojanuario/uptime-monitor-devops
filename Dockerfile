FROM golang:1.26.5-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o uptime-monitor ./cmd/api

FROM alpine:3.24.1

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/uptime-monitor .

EXPOSE 8080

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["./uptime-monitor"]
