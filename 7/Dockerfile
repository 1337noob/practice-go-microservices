FROM golang:1.26.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/user ./services/user/cmd/main.go && \
    go build -o bin/notification ./services/notification/cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bin .
EXPOSE 50051
CMD ["./user"]
