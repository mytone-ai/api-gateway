FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
