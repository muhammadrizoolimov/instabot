package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"instabot2/internal/downloader"
	"instabot2/internal/services"
	"instabot2/internal/worker"
	"log"
	"strings"
)

// HandleMessage processes incoming messages
func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	if message.Text == "" {
		return
	}

	text := strings.TrimSpace(message.Text)

	// Handle commands
	if strings.HasPrefix(text, "/") {
		b.handleCommand(message, text)
		return
	}

	// Check if message contains a URL
	if strings.Contains(text, "http://") || strings.Contains(text, "https://") {
		b.handleURL(message, text)
		return
	}

	// Default response
	b.sendMessage(message.Chat.ID, "❓ Iltimos, video havolasini yuboring yoki /help komandasidan foydalaning.")
}

// handleCommand processes bot commands
func (b *Bot) handleCommand(message *tgbotapi.Message, command string) {
	chatID := message.Chat.ID

	switch {
	case strings.HasPrefix(command, "/start"):
		welcomeText := `👋 <b>Assalomu alaykum!</b>

Men Instagram, TikTok, YouTube, Pinterest, Snapchat va Likee'dan video yuklab beraman.

<b>Qanday foydalanish:</b>
1️⃣ Video havolasini yuboring
2️⃣ Video ostidagi tugmalardan foydalaning:
   🎵 <b>Musiqasini topish</b> - Video musiqasini qidirish
   ➕ <b>Guruhga qo'shish</b> - Botni guruhga qo'shish

<b>Qo'llab-quvvatlanadigan platformalar:</b>
• Instagram (post, reel, story)
• TikTok
• YouTube (video, shorts)
• Pinterest
• Snapchat
• Likee

Havolani yuboring va men videoni yuklab beraman! 🚀`

		msg := tgbotapi.NewMessage(chatID, welcomeText)
		msg.ParseMode = "HTML"
		b.API.Send(msg)

	case strings.HasPrefix(command, "/help"):
		helpText := `📚 <b>Yordam</b>

<b>Asosiy funksiyalar:</b>

1️⃣ <b>Video yuklash:</b>
   - Video havolasini yuboring
   - Bot videoni yuklaydi va yuboradi

2️⃣ <b>Musiqa qidirish:</b>
   - Video ostidagi "🎵 Musiqasini topish" tugmasini bosing
   - 10 ta natija ko'rsatiladi
   - Kerakli musiqani tanlang va yuklang

3️⃣ <b>Guruhda ishlash:</b>
   - "➕ Guruhga qo'shish" tugmasini bosing
   - Botni guruhga admin qiling
   - Guruhda havola yuboring va bot javob beradi

<b>Qo'llab-quvvatlanadigan platformalar:</b>
Instagram, TikTok, YouTube, Pinterest, Snapchat, Likee

Savol bo'lsa, @support ga murojaat qiling.`

		msg := tgbotapi.NewMessage(chatID, helpText)
		msg.ParseMode = "HTML"
		b.API.Send(msg)

	default:
		b.sendMessage(chatID, "❓ Noma'lum komanda. /help ni ko'ring.")
	}
}

// handleURL processes URLs from supported platforms
func (b *Bot) handleURL(message *tgbotapi.Message, text string) {
	chatID := message.Chat.ID

	// Extract URL from text
	words := strings.Fields(text)
	var url string
	for _, word := range words {
		if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
			url = word
			break
		}
	}

	if url == "" {
		b.sendMessage(chatID, "❌ Havola topilmadi.")
		return
	}

	// Check if platform is supported
	if !downloader.IsSupportedURL(url) {
		b.sendMessage(chatID, "❌ Ushbu platforma qo'llab-quvvatlanmaydi.\n\nQo'llab-quvvatlanadigan: Instagram, TikTok, YouTube, Pinterest, Snapchat, Likee")
		return
	}

	platform := downloader.DetectPlatform(url)
	log.Printf("Downloading from %s: %s", platform, url)

	// Show loading indicator
	loadingMsg := b.sendMessage(chatID, "⏳ Video yuklanmoqda...")

	// Submit to worker pool
	pool := worker.GetPool()
	pool.Submit(worker.Job{
		ID:  fmt.Sprintf("video_%d_%s", chatID, url),
		URL: url,
		Handler: func(videoURL string) error {
			return b.downloadAndSendVideo(chatID, loadingMsg.MessageID, videoURL, platform)
		},
	})
}

// downloadAndSendVideo downloads video and sends to user with smart buttons
func (b *Bot) downloadAndSendVideo(chatID int64, loadingMsgID int, url string, platform string) error {
	defer b.deleteMessage(chatID, loadingMsgID)

	// Download video
	mediaService := services.NewMediaService(b.TempDir)
	result, err := mediaService.DownloadVideo(url)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Yuklash xatosi: %v", err))
		return err
	}
	defer mediaService.Cleanup(result.FilePath)

	// Create smart buttons
	keyboard := b.createSmartButtons(result.Title)

	// Send video with buttons
	video := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(result.FilePath))
	video.Caption = fmt.Sprintf("✅ %s\n\n📱 Platform: %s", result.Title, platform)
	video.ReplyMarkup = keyboard

	_, err = b.API.Send(video)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Yuborish xatosi: %v", err))
		return err
	}

	return nil
}

// createSmartButtons creates inline keyboard with smart buttons
func (b *Bot) createSmartButtons(videoTitle string) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("🎵 Musiqasini topish", fmt.Sprintf("search_music:%s", videoTitle)),
		},
		{
			tgbotapi.NewInlineKeyboardButtonURL("➕ Guruhga qo'shish", "https://t.me/saveinstabot?startgroup=true"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// sendMessage sends a text message and returns the sent message
func (b *Bot) sendMessage(chatID int64, text string) *tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	sent, err := b.API.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return nil
	}
	return &sent
}

// deleteMessage deletes a message
func (b *Bot) deleteMessage(chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	b.API.Request(deleteMsg)
}
