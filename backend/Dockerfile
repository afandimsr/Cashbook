# backend/Dockerfile

# Stage 1: Build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary (CGO_ENABLED=0 agar bisa jalan di alpine/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Copy .env (Opsional: docker-compose bisa handle ini, tapi aman jika dicopy)
COPY .env .

EXPOSE 8181

CMD ["./main"]