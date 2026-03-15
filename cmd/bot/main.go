package main

import (
	"instabot2/internal/bot"
	"instabot2/internal/database"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Get configuration from environment
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/bot.db"
	}

	tempDir := os.Getenv("TEMP_DIR")
	if tempDir == "" {
		tempDir = "./temp"
	}

	// Initialize database
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database initialized successfully")

	// Create bot instance
	b, err := bot.New(token, db, tempDir)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down bot...")
		os.Exit(0)
	}()

	// Start bot
	log.Println("Starting InstaBot2...")
	if err := b.Start(); err != nil {
		log.Fatalf("Bot error: %v", err)
	}
}
