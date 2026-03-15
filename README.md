# InstaBot2 - Multi-Platform Video Downloader Bot

Telegram bot for downloading videos from Instagram, TikTok, YouTube, Pinterest, Snapchat, and Likee with integrated music search functionality.

## Features

### 🎥 Video Download
- **Supported Platforms**: Instagram, TikTok, YouTube, Pinterest, Snapchat, Likee
- **Smart Buttons**: Every downloaded video includes:
  - 🎵 **Musiqasini topish** - Search for music based on video title
  - ➕ **Guruhga qo'shish** - Add bot to group

### 🎵 Music Search Module
- **yt-dlp Integration**: Search music on YouTube
- **Pagination**: 10 results per page with ⬅️❌➡️ navigation
- **UTF-8 Cleaning**: Safe character handling for Telegram
- **HTML Formatting**: Clean and readable results

### 🔊 Audio Extraction
- Extract audio from any supported video URL
- MP3 format with best quality
- Automatic metadata extraction (title, duration, performer)

### ⚡ Worker Pool System
- **maxWorkers=5**: Concurrent download processing
- **Queue Management**: Prevents duplicate jobs
- **Non-blocking**: Downloads don't freeze the bot

### 💾 Caching System
- **music_cache Table**: Stores downloaded audio file_id
- **Fast Re-sends**: Cached files sent instantly without re-download
- **SQLite Database**: Lightweight and efficient

### 🎨 UX Improvements
- **Loading Indicators**: ⏳ shows during processing
- **Auto-cleanup**: Loading messages deleted when done
- **Safe Filenames**: Special characters filtered
- **Error Handling**: User-friendly error messages in Uzbek

## Installation

### Prerequisites
- Go 1.21+
- yt-dlp installed and in PATH
- Telegram Bot Token

### Setup

1. **Clone the repository:**
```bash
git clone <repo-url>
cd instabot
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Install yt-dlp:**
```bash
# Linux/macOS
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

# Or using pip
pip install yt-dlp
```

4. **Configure environment variables:**
```bash
cp .env.example .env
# Edit .env and add your TELEGRAM_BOT_TOKEN
```

5. **Run the bot:**
```bash
export TELEGRAM_BOT_TOKEN="your_bot_token_here"
go run cmd/bot/main.go
```

## Project Structure

```
instabot2/
├── cmd/
│   └── bot/
│       └── main.go              # Entry point
├── internal/
│   ├── bot/
│   │   ├── bot.go               # Bot core
│   │   ├── handlers.go          # Message handlers
│   │   └── callback.go          # Callback query handlers
│   ├── music/
│   │   └── search.go            # Music search with yt-dlp
│   ├── services/
│   │   └── media.go             # Video/Audio download & extraction
│   ├── worker/
│   │   └── worker.go            # Worker pool (maxWorkers=5)
│   ├── database/
│   │   └── db.go                # SQLite with music_cache
│   ├── downloader/
│   │   └── platforms.go         # Platform detection
│   └── utils/
│       └── sanitize.go          # UTF-8 & HTML sanitization
├── config/
├── go.mod
└── README.md
```

## Usage

### Download Video
1. Send a video URL from any supported platform
2. Bot downloads and sends the video
3. Use the inline buttons below the video:
   - **🎵 Musiqasini topish**: Search for similar music
   - **➕ Guruhga qo'shish**: Add bot to a group

### Search Music
1. After downloading a video, click "🎵 Musiqasini topish"
2. Bot searches for music based on the video title
3. Browse results with ⬅️➡️ pagination (10 per page)
4. Click on a result number to download the audio

### Commands
- `/start` - Welcome message and instructions
- `/help` - Detailed help information

## Technical Details

### Music Search
- **Engine**: yt-dlp with YouTube search
- **Results**: 10 per page with pagination
- **Format**: JSON parsing with UTF-8 safety
- **Buttons**: Dynamic inline keyboard (1-10)

### Audio Extraction
- **Format**: MP3 (best quality)
- **Metadata**: Title, duration, performer
- **Storage**: File cached by file_id in database

### Worker Pool
- **Concurrency**: 5 workers maximum
- **Queue**: Channel-based job distribution
- **Safety**: Prevents duplicate processing

### Database Schema
```sql
CREATE TABLE music_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    title TEXT NOT NULL,
    file_id TEXT NOT NULL,
    duration INTEGER,
    performer TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(query, title)
);
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEGRAM_BOT_TOKEN` | Telegram Bot API Token | *Required* |
| `DATABASE_PATH` | SQLite database file path | `./data/bot.db` |
| `TEMP_DIR` | Temporary files directory | `./temp` |

## Dependencies

- **telegram-bot-api/v5**: Telegram Bot API wrapper
- **mattn/go-sqlite3**: SQLite3 driver
- **yt-dlp**: External tool for downloading

## License

MIT License

## Support

For questions or issues, contact @support on Telegram.
