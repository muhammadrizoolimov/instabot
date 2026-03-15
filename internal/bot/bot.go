package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"instabot2/internal/database"
	"log"
)

type Bot struct {
	API     *tgbotapi.BotAPI
	DB      *database.Database
	TempDir string
}

// New creates a new bot instance
func New(token string, db *database.Database, tempDir string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		API:     api,
		DB:      db,
		TempDir: tempDir,
	}, nil
}

// Start starts the bot and listens for updates
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	log.Println("Bot started. Waiting for updates...")

	for update := range updates {
		go b.handleUpdate(update)
	}

	return nil
}

// handleUpdate routes updates to appropriate handlers
func (b *Bot) handleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	if update.Message != nil {
		b.HandleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.HandleCallback(update.CallbackQuery)
	}
}
