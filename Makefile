.PHONY: build run clean test deps install

# Build the bot
build:
	go build -o bin/instabot2 cmd/bot/main.go

# Run the bot
run:
	go run cmd/bot/main.go

# Clean build artifacts
clean:
	rm -rf bin/ temp/ data/

# Run tests
test:
	go test -v ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install yt-dlp (Linux/macOS)
install-ytdlp:
	@echo "Installing yt-dlp..."
	@curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /tmp/yt-dlp
	@sudo mv /tmp/yt-dlp /usr/local/bin/yt-dlp
	@sudo chmod a+rx /usr/local/bin/yt-dlp
	@echo "yt-dlp installed successfully!"

# Full setup
setup: deps install-ytdlp
	@echo "Setup complete! Configure .env and run 'make run'"
