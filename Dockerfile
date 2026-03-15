FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o /app/bin/instabot2 cmd/bot/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    python3 \
    py3-pip \
    sqlite

# Install yt-dlp
RUN pip3 install --no-cache-dir yt-dlp

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/instabot2 /app/instabot2

# Create directories
RUN mkdir -p /app/data /app/temp

# Environment variables
ENV DATABASE_PATH=/app/data/bot.db
ENV TEMP_DIR=/app/temp

# Run the bot
CMD ["/app/instabot2"]
