FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . . --exclude tmp .git 
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/server/main.go

FROM debian:stable-slim
RUN apt-get update && apt-get install -y smartmontools iproute2 dmidecode && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /api /api
EXPOSE 8080
CMD ["/api"]
