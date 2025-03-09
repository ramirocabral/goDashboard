# frontend build
FROM node:18 AS frontend-builder
WORKDIR /client

COPY client/package.json client/package-lock.json ./
RUN npm install

COPY client ./
RUN npm run build

# backend build
FROM golang:1.23 AS backend-builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/server/main.go

# final stage
FROM debian:stable-slim
WORKDIR /app

RUN apt-get update && apt-get install -y smartmontools iproute2 dmidecode && \
    rm -rf /var/lib/apt/lists/*

COPY --from=backend-builder /api /api

COPY --from=frontend-builder /client/dist /dist

EXPOSE 8080 

CMD ["/api"]
