FROM golang:1.25.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/pr

FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/app .

COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./app"]