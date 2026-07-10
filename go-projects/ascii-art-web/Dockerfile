# ───── Stage 1: Build ─────
FROM golang:1.23.2 AS builder

LABEL maintainer="Manar Mohamed"
LABEL version="1.0"
LABEL description="ASCII Art Web App - Dockerized using multi-stage build"

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -o server main.go

# ───── Stage 2: Lightweight Final Image ─────
FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

# Copy only what's needed
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/banners ./banners

EXPOSE 8080

CMD [ "./server" ]