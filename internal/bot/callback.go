package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"instabot2/internal/music"
	"instabot2/internal/services"
	"instabot2/internal/utils"
	"instabot2/internal/worker"
	"log"
	"strconv"
	"strings"
)

// HandleCallback processes callback queries from inline buttons
func (b *Bot) HandleCallback(callback *tgbotapi.CallbackQuery) {
	// Answer callback to remove loading state
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := b.API.Request(callbackConfig); err != nil {
		log.Printf("Error answering callback: %v", err)
	}

	data := callback.Data

	switch {
	case strings.HasPrefix(data, "search_music:"):
		// Extract video title from callback data
		parts := strings.SplitN(data, ":", 2)
		if len(parts) != 2 {
			b.sendMessage(callback.Message.Chat.ID, "❌ Xatolik: noto'g'ri ma'lumot")
			return
		}
		videoTitle := parts[1]
		b.handleMusicSearch(callback.Message, videoTitle, 1)

	case strings.HasPrefix(data, "music_page:"):
		// Pagination for music search
		parts := strings.Split(data, ":")
		if len(parts) != 3 {
			return
		}
		query := parts[1]
		page, _ := strconv.Atoi(parts[2])
		b.handleMusicSearch(callback.Message, query, page)

	case strings.HasPrefix(data, "download_music:"):
		// Download selected music
		parts := strings.SplitN(data, ":", 2)
		if len(parts) != 2 {
			return
		}
		musicURL := parts[1]
		b.handleMusicDownload(callback.Message, musicURL)

	case data == "close_music":
		// Close music search results
		deleteMsg := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		b.API.Request(deleteMsg)
	}
}

// handleMusicSearch searches for music and displays results with pagination
func (b *Bot) handleMusicSearch(message *tgbotapi.Message, query string, page int) {
	chatID := message.Chat.ID

	// Show loading indicator
	loadingMsg := b.sendMessage(chatID, "⏳ Musiqa qidirilmoqda...")

	// Search for music
	results, err := music.SearchMusic(query, page)
	if err != nil {
		b.deleteMessage(chatID, loadingMsg.MessageID)
		b.sendMessage(chatID, fmt.Sprintf("❌ Qidiruv xatosi: %v", err))
		return
	}

	// Delete loading message
	b.deleteMessage(chatID, loadingMsg.MessageID)

	// Build results message
	text := fmt.Sprintf("🎵 <b>Qidiruv natijalari:</b> %s\n\n", utils.FormatHTML(query))

	startNum := (page-1)*music.ResultsPerPage + 1
	for i, result := range results.Results {
		num := startNum + i
		duration := music.FormatDuration(result.Duration)
		title := utils.TruncateString(result.Title, 100)
		text += fmt.Sprintf("%d. <b>%s</b>\n⏱ %s\n\n", num, utils.FormatHTML(title), duration)
	}

	text += fmt.Sprintf("📄 Sahifa: %d/%d", page, results.TotalPages)

	// Create inline keyboard with music results and pagination
	var buttons [][]tgbotapi.InlineKeyboardButton

	// Music download buttons (max 10)
	for i, result := range results.Results {
		num := startNum + i
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d - Yuklash", num),
			fmt.Sprintf("download_music:%s", result.URL),
		)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{btn})
	}

	// Pagination row
	var paginationRow []tgbotapi.InlineKeyboardButton
	if page > 1 {
		paginationRow = append(paginationRow,
			tgbotapi.NewInlineKeyboardButtonData("⬅️", fmt.Sprintf("music_page:%s:%d", query, page-1)),
		)
	}
	paginationRow = append(paginationRow,
		tgbotapi.NewInlineKeyboardButtonData("❌", "close_music"),
	)
	if page < results.TotalPages {
		paginationRow = append(paginationRow,
			tgbotapi.NewInlineKeyboardButtonData("➡️", fmt.Sprintf("music_page:%s:%d", query, page+1)),
		)
	}
	buttons = append(buttons, paginationRow)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	// Update or send new message
	if message.ReplyMarkup != nil {
		// Edit existing message
		edit := tgbotapi.NewEditMessageText(chatID, message.MessageID, text)
		edit.ParseMode = "HTML"
		edit.ReplyMarkup = &keyboard
		b.API.Send(edit)
	} else {
		// Send new message
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "HTML"
		msg.ReplyMarkup = keyboard
		b.API.Send(msg)
	}
}

// handleMusicDownload downloads music using worker pool
func (b *Bot) handleMusicDownload(message *tgbotapi.Message, url string) {
	chatID := message.Chat.ID

	// Show loading indicator
	loadingMsg := b.sendMessage(chatID, "⏳ Musiqa yuklanmoqda...")

	// Submit to worker pool
	pool := worker.GetPool()
	pool.Submit(worker.Job{
		ID:  fmt.Sprintf("music_%d_%s", chatID, url),
		URL: url,
		Handler: func(musicURL string) error {
			return b.downloadAndSendMusic(chatID, loadingMsg.MessageID, musicURL)
		},
	})
}

// downloadAndSendMusic downloads music and sends to user
func (b *Bot) downloadAndSendMusic(chatID int64, loadingMsgID int, url string) error {
	defer b.deleteMessage(chatID, loadingMsgID)

	// Check cache first
	cachedFileID, err := b.DB.GetCachedMusic(url, url)
	if err == nil && cachedFileID != "" {
		audio := tgbotapi.NewAudio(chatID, tgbotapi.FileID(cachedFileID))
		b.API.Send(audio)
		return nil
	}

	// Download audio
	mediaService := services.NewMediaService(b.TempDir)
	result, err := mediaService.ExtractAudio(url)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Yuklash xatosi: %v", err))
		return err
	}
	defer mediaService.Cleanup(result.FilePath)

	// Send audio file
	audio := tgbotapi.NewAudio(chatID, tgbotapi.FilePath(result.FilePath))
	audio.Title = result.Title
	audio.Duration = result.Duration

	sent, err := b.API.Send(audio)
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("❌ Yuborish xatosi: %v", err))
		return err
	}

	// Cache file_id if available
	if len(sent.Audio.FileID) > 0 {
		b.DB.CacheMusic(url, result.Title, sent.Audio.FileID, result.Duration, "")
	}

	return nil
}
