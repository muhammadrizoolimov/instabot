# API Documentation

## Internal APIs

### Music Search API

#### `SearchMusic(query string, page int) (*SearchResponse, error)`

Searches for music using yt-dlp.

**Parameters:**
- `query`: Search query (auto-sanitized for UTF-8)
- `page`: Page number (1-based indexing)

**Returns:**
```go
type SearchResponse struct {
    Results    []SearchResult  // 10 results per page
    Page       int             // Current page number
    TotalPages int             // Total pages available
}

type SearchResult struct {
    Title    string  // Sanitized track title
    URL      string  // YouTube URL
    Duration int     // Duration in seconds
    Uploader string  // Artist/channel name
}
```

**Example:**
```go
results, err := music.SearchMusic("Uzbek music", 1)
if err != nil {
    log.Fatal(err)
}

for i, result := range results.Results {
    fmt.Printf("%d. %s - %s\n", i+1, result.Title, music.FormatDuration(result.Duration))
}
```

---

### Media Service API

#### `DownloadVideo(url string) (*DownloadResult, error)`

Downloads video from supported platform.

**Parameters:**
- `url`: Video URL from supported platform

**Returns:**
```go
type DownloadResult struct {
    FilePath  string  // Local file path
    Title     string  // Video title (sanitized)
    Thumbnail string  // Thumbnail URL (if available)
    Duration  int     // Duration in seconds
}
```

**Example:**
```go
ms := services.NewMediaService("./temp")
result, err := ms.DownloadVideo("https://www.youtube.com/watch?v=...")
if err != nil {
    log.Fatal(err)
}
defer ms.Cleanup(result.FilePath)

// Use result.FilePath to send video
```

#### `ExtractAudio(url string) (*DownloadResult, error)`

Extracts audio from video URL.

**Parameters:**
- `url`: Video/audio URL

**Returns:**
- Same as `DownloadVideo` but FilePath points to MP3 file

**Example:**
```go
ms := services.NewMediaService("./temp")
audio, err := ms.ExtractAudio("https://www.youtube.com/watch?v=...")
if err != nil {
    log.Fatal(err)
}
defer ms.Cleanup(audio.FilePath)

// Send audio.FilePath as Telegram audio
```

---

### Worker Pool API

#### `GetPool() *WorkerPool`

Returns singleton worker pool instance.

#### `Submit(job Job)`

Submits job to worker pool queue.

**Job Structure:**
```go
type Job struct {
    ID      string                  // Unique job identifier
    URL     string                  // URL to process
    Handler func(string) error      // Job handler function
}
```

**Example:**
```go
pool := worker.GetPool()
pool.Submit(worker.Job{
    ID:  "download_12345",
    URL: "https://...",
    Handler: func(url string) error {
        // Download and process
        return nil
    },
})
```

#### `Wait()`

Waits for all jobs to complete.

```go
pool.Wait()
```

---

### Database API

#### `New(dbPath string) (*Database, error)`

Creates new database connection and initializes tables.

**Example:**
```go
db, err := database.New("./data/bot.db")
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

#### `CacheMusic(query, title, fileID string, duration int, performer string) error`

Stores music cache entry.

**Parameters:**
- `query`: Original search query
- `title`: Track title
- `fileID`: Telegram file_id
- `duration`: Duration in seconds
- `performer`: Artist name

**Example:**
```go
err := db.CacheMusic(
    "uzbek music",
    "Beautiful Song",
    "AgACAgIAAxkBAAI...",
    180,
    "Artist Name",
)
```

#### `GetCachedMusic(query, title string) (string, error)`

Retrieves cached music file_id.

**Returns:**
- `string`: Telegram file_id
- `error`: sql.ErrNoRows if not found

**Example:**
```go
fileID, err := db.GetCachedMusic("uzbek music", "Beautiful Song")
if err == nil {
    // Use cached fileID
    audio := tgbotapi.NewAudio(chatID, tgbotapi.FileID(fileID))
    bot.Send(audio)
}
```

---

### Utilities API

#### `SanitizeForTelegram(s string) string`

Removes/replaces characters that cause Telegram API errors.

**Features:**
- Unescapes HTML entities
- Filters non-alphanumeric characters
- Keeps safe punctuation
- Cleans multiple spaces
- Truncates to 200 characters

**Example:**
```go
dirty := "Song with <weird> chars & émojis 🎵"
clean := utils.SanitizeForTelegram(dirty)
// Result: "Song with weird chars & emojis"
```

#### `FormatHTML(text string) string`

Escapes special HTML characters for Telegram HTML parsing.

**Example:**
```go
html := utils.FormatHTML("<b>Bold</b> & <i>italic</i>")
// Result: "&lt;b&gt;Bold&lt;/b&gt; &amp; &lt;i&gt;italic&lt;/i&gt;"
```

#### `TruncateString(s string, maxLen int) string`

Truncates string to maxLen and adds ellipsis.

**Example:**
```go
long := "Very long string that needs truncation"
short := utils.TruncateString(long, 20)
// Result: "Very long string tha..."
```

---

### Platform Detection API

#### `DetectPlatform(url string) string`

Detects platform from URL.

**Returns:**
- Platform name: "Instagram", "TikTok", "YouTube", etc.
- "Unknown" if not supported

**Example:**
```go
platform := downloader.DetectPlatform("https://www.instagram.com/p/...")
// Result: "Instagram"
```

#### `IsSupportedURL(url string) bool`

Checks if URL is from supported platform.

**Example:**
```go
if downloader.IsSupportedURL(url) {
    // Process download
}
```

---

## Telegram Callback Data Format

### Music Search
```
search_music:<video_title>
```
Triggers music search based on video title.

### Music Pagination
```
music_page:<query>:<page_number>
```
Navigate to specific page of search results.

**Example:**
```
music_page:uzbek music:2
```

### Music Download
```
download_music:<youtube_url>
```
Downloads audio from YouTube URL.

**Example:**
```
download_music:https://www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Close Results
```
close_music
```
Deletes music search results message.

---

## Message Flow Examples

### Complete Video Download Flow

```go
// 1. User sends URL
update.Message.Text = "https://www.youtube.com/watch?v=..."

// 2. Bot detects platform
platform := downloader.DetectPlatform(url)

// 3. Submit to worker
pool.Submit(worker.Job{
    ID: "video_123456",
    URL: url,
    Handler: func(url string) error {
        // 4. Download video
        ms := services.NewMediaService("./temp")
        result, _ := ms.DownloadVideo(url)
        
        // 5. Create smart buttons
        keyboard := createSmartButtons(result.Title)
        
        // 6. Send video
        video := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(result.FilePath))
        video.ReplyMarkup = keyboard
        bot.Send(video)
        
        // 7. Cleanup
        ms.Cleanup(result.FilePath)
        return nil
    },
})
```

### Complete Music Search Flow

```go
// 1. User clicks "🎵 Musiqasini topish"
callback.Data = "search_music:Beautiful Video Title"

// 2. Extract title
parts := strings.SplitN(callback.Data, ":", 2)
title := parts[1]

// 3. Search music
results, _ := music.SearchMusic(title, 1)

// 4. Build keyboard with results
var buttons [][]tgbotapi.InlineKeyboardButton
for i, result := range results.Results {
    btn := tgbotapi.NewInlineKeyboardButtonData(
        fmt.Sprintf("%d - Yuklash", i+1),
        fmt.Sprintf("download_music:%s", result.URL),
    )
    buttons = append(buttons, []tgbotapi.InlineKeyboardButton{btn})
}

// 5. Add pagination
paginationRow := []tgbotapi.InlineKeyboardButton{
    tgbotapi.NewInlineKeyboardButtonData("⬅️", "music_page:"+title+":0"),
    tgbotapi.NewInlineKeyboardButtonData("❌", "close_music"),
    tgbotapi.NewInlineKeyboardButtonData("➡️", "music_page:"+title+":2"),
}
buttons = append(buttons, paginationRow)

// 6. Send results
keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)
msg := tgbotapi.NewMessage(chatID, resultsText)
msg.ReplyMarkup = keyboard
bot.Send(msg)
```

### Complete Music Download with Caching

```go
// 1. User clicks "1 - Yuklash"
callback.Data = "download_music:https://www.youtube.com/watch?v=..."

// 2. Check cache
fileID, err := db.GetCachedMusic(url, "")
if err == nil {
    // Cache hit - instant send
    audio := tgbotapi.NewAudio(chatID, tgbotapi.FileID(fileID))
    bot.Send(audio)
    return
}

// 3. Cache miss - download
ms := services.NewMediaService("./temp")
result, _ := ms.ExtractAudio(url)

// 4. Send audio
audio := tgbotapi.NewAudio(chatID, tgbotapi.FilePath(result.FilePath))
audio.Title = result.Title
audio.Duration = result.Duration
sent, _ := bot.Send(audio)

// 5. Cache for future
if sent.Audio != nil {
    db.CacheMusic(url, result.Title, sent.Audio.FileID, result.Duration, "")
}

// 6. Cleanup
ms.Cleanup(result.FilePath)
```

---

## Error Handling

All APIs return errors following Go conventions. Always check errors:

```go
result, err := music.SearchMusic(query, 1)
if err != nil {
    // Handle error
    log.Printf("Search failed: %v", err)
    return
}
```

Common errors:
- `yt-dlp search failed`: yt-dlp not installed or network error
- `no results found`: Empty search results
- `audio extraction failed`: Invalid URL or unsupported format
- `sql.ErrNoRows`: Cache miss (not an error, just not cached)
