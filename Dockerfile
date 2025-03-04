FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/server/main.go

FROM debian:buster-slim
RUN apt-get update && apt-get install -y smartmontools iproute2 && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /api /api
EXPOSE 8080
CMD ["/api"]
