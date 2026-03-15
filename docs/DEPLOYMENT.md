# Deployment Guide

## Prerequisites

- Linux server (Ubuntu 20.04+ recommended)
- Docker & Docker Compose (recommended) OR Go 1.21+
- yt-dlp installed
- Telegram Bot Token from [@BotFather](https://t.me/BotFather)

---

## Method 1: Docker Deployment (Recommended)

### Step 1: Clone Repository

```bash
git clone <repository-url>
cd instabot
```

### Step 2: Configure Environment

```bash
cp .env.example .env
nano .env
```

Add your bot token:
```env
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
DATABASE_PATH=./data/bot.db
```

### Step 3: Build and Run

```bash
docker-compose up -d
```

### Step 4: Check Logs

```bash
docker-compose logs -f instabot2
```

You should see:
```
Database initialized successfully
Authorized on account YourBotName
Bot started. Waiting for updates...
```

### Step 5: Test Bot

Open Telegram and send `/start` to your bot.

### Management Commands

```bash
# Stop bot
docker-compose stop

# Restart bot
docker-compose restart

# View logs
docker-compose logs -f

# Rebuild after code changes
docker-compose up -d --build

# Stop and remove
docker-compose down
```

---

## Method 2: Manual Deployment

### Step 1: Install Dependencies

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install yt-dlp
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

# Install ffmpeg (required by yt-dlp for audio extraction)
sudo apt install -y ffmpeg

# Verify installations
go version
yt-dlp --version
ffmpeg -version
```

### Step 2: Clone and Build

```bash
git clone <repository-url>
cd instabot
go mod download
go build -o bin/instabot2 cmd/bot/main.go
```

### Step 3: Configure

```bash
cp .env.example .env
nano .env
```

### Step 4: Run Bot

```bash
export TELEGRAM_BOT_TOKEN="your_token_here"
./bin/instabot2
```

### Step 5: Run as Service (systemd)

Create service file:

```bash
sudo nano /etc/systemd/system/instabot2.service
```

Add content:
```ini
[Unit]
Description=InstaBot2 Telegram Bot
After=network.target

[Service]
Type=simple
User=your_username
WorkingDirectory=/path/to/instabot
Environment="TELEGRAM_BOT_TOKEN=your_token_here"
Environment="DATABASE_PATH=/path/to/instabot/data/bot.db"
Environment="TEMP_DIR=/path/to/instabot/temp"
ExecStart=/path/to/instabot/bin/instabot2
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable instabot2
sudo systemctl start instabot2
sudo systemctl status instabot2
```

View logs:
```bash
sudo journalctl -u instabot2 -f
```

---

## Method 3: Using Makefile

### Quick Start

```bash
# Install dependencies and yt-dlp
make setup

# Build binary
make build

# Run bot
export TELEGRAM_BOT_TOKEN="your_token"
make run
```

### Available Commands

```bash
make build         # Build binary to bin/instabot2
make run           # Run bot directly
make clean         # Clean build artifacts
make test          # Run tests
make deps          # Download Go dependencies
make install-ytdlp # Install yt-dlp
make setup         # Full setup (deps + yt-dlp)
```

---

## Configuration Options

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TELEGRAM_BOT_TOKEN` | Bot token from @BotFather | - | Yes |
| `DATABASE_PATH` | SQLite database file path | `./data/bot.db` | No |
| `TEMP_DIR` | Temporary download directory | `./temp` | No |

### Database Location

By default, database is stored in `./data/bot.db`. For production:

```bash
mkdir -p /var/lib/instabot2
export DATABASE_PATH=/var/lib/instabot2/bot.db
```

### Temp Directory

Temporary files are stored in `./temp` and auto-deleted after use. Ensure sufficient disk space:

```bash
mkdir -p /tmp/instabot2
export TEMP_DIR=/tmp/instabot2
```

---

## Monitoring

### Health Checks

Check if bot is responsive:

```bash
# Docker
docker-compose logs --tail=50 instabot2

# Systemd
sudo journalctl -u instabot2 --since "10 minutes ago"
```

### Database Size

Monitor database growth:

```bash
du -h data/bot.db
sqlite3 data/bot.db "SELECT COUNT(*) FROM music_cache;"
```

### Disk Space

Monitor temp directory:

```bash
du -sh temp/
```

Set up auto-cleanup (cron):

```bash
# Add to crontab
0 */6 * * * find /path/to/instabot/temp -type f -mtime +1 -delete
```

---

## Scaling

### Single Server Optimization

1. **Increase Workers:**
   Edit `internal/worker/worker.go`:
   ```go
   const maxWorkers = 10  // Increase from 5
   ```

2. **Increase Queue Size:**
   ```go
   jobs: make(chan Job, 200),  // Increase from 100
   ```

3. **Use PostgreSQL:**
   Replace SQLite with PostgreSQL for better concurrency.

### Multi-Server Deployment

1. **Load Balancer:**
   - Deploy multiple bot instances
   - Use Telegram webhooks with load balancer
   - Share database (PostgreSQL)

2. **Shared Cache:**
   - Use Redis for file_id caching
   - Multiple instances can share cache

3. **Distributed Workers:**
   - Use RabbitMQ/Redis for job queue
   - Separate worker servers for downloads

---

## Backup & Recovery

### Backup Database

```bash
# Create backup
sqlite3 data/bot.db ".backup 'data/bot.db.backup'"

# Or use cp
cp data/bot.db data/bot.db.$(date +%Y%m%d)
```

### Automated Backups

Add to crontab:
```bash
0 2 * * * sqlite3 /path/to/instabot/data/bot.db ".backup '/backup/bot.db.$(date +\%Y\%m\%d)'"
```

### Restore

```bash
cp data/bot.db.backup data/bot.db
# Restart bot
sudo systemctl restart instabot2
```

---

## Troubleshooting

### Bot Not Responding

1. **Check bot is running:**
   ```bash
   # Docker
   docker-compose ps
   
   # Systemd
   sudo systemctl status instabot2
   ```

2. **Check logs:**
   ```bash
   # Docker
   docker-compose logs --tail=100 instabot2
   
   # Systemd
   sudo journalctl -u instabot2 -n 100
   ```

3. **Verify token:**
   ```bash
   curl https://api.telegram.org/bot<TOKEN>/getMe
   ```

### yt-dlp Errors

1. **Update yt-dlp:**
   ```bash
   sudo yt-dlp -U
   ```

2. **Test manually:**
   ```bash
   yt-dlp --dump-json "ytsearch:test"
   ```

3. **Check ffmpeg:**
   ```bash
   ffmpeg -version
   ```

### Download Failures

1. **Check disk space:**
   ```bash
   df -h
   ```

2. **Check temp directory permissions:**
   ```bash
   ls -la temp/
   chmod 755 temp/
   ```

3. **Verify yt-dlp can access URL:**
   ```bash
   yt-dlp -F "https://youtube.com/watch?v=..."
   ```

### Database Locked

SQLite may lock under heavy load:

1. **Increase timeout:**
   Add to `database.New()`:
   ```go
   db.Exec("PRAGMA busy_timeout = 5000")
   ```

2. **Switch to PostgreSQL** for production.

---

## Security Hardening

### 1. File Permissions

```bash
chmod 600 .env
chmod 700 data/
chmod 755 temp/
```

### 2. Firewall (UFW)

```bash
sudo ufw allow ssh
sudo ufw enable
# No need to open ports for polling bot
```

### 3. User Isolation

```bash
sudo useradd -r -s /bin/false instabot
sudo chown -R instabot:instabot /path/to/instabot
```

Update systemd service:
```ini
User=instabot
Group=instabot
```

### 4. Rate Limiting

Implement in code or use fail2ban to block abusive users.

---

## Updates & Maintenance

### Update Bot Code

```bash
# Docker
git pull
docker-compose up -d --build

# Manual
git pull
go build -o bin/instabot2 cmd/bot/main.go
sudo systemctl restart instabot2
```

### Update Dependencies

```bash
go get -u ./...
go mod tidy
```

### Update yt-dlp

```bash
sudo yt-dlp -U
```

---

## Performance Tuning

### 1. Database Optimization

```sql
-- Run periodically
VACUUM;
ANALYZE;

-- Add more indexes if needed
CREATE INDEX idx_music_cache_created ON music_cache(created_at);
```

### 2. Temp File Cleanup

```bash
# Add to cleanup.sh
#!/bin/bash
find /path/to/temp -type f -mtime +1 -delete

# Run via cron
0 */6 * * * /path/to/cleanup.sh
```

### 3. Log Rotation

```bash
# For systemd logs
sudo journalctl --vacuum-time=7d
```

---

## Monitoring & Alerts

### Basic Monitoring Script

```bash
#!/bin/bash
# monitor.sh

# Check if bot is running
if ! systemctl is-active --quiet instabot2; then
    echo "Bot is down!" | mail -s "InstaBot2 Alert" admin@example.com
fi

# Check disk space
USAGE=$(df -h /path/to/instabot | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $USAGE -gt 80 ]; then
    echo "Disk usage: ${USAGE}%" | mail -s "Disk Alert" admin@example.com
fi
```

Run via cron:
```bash
*/5 * * * * /path/to/monitor.sh
```

---

## Production Checklist

- [ ] Bot token configured securely
- [ ] Database path set to persistent location
- [ ] yt-dlp installed and updated
- [ ] ffmpeg installed
- [ ] Systemd service configured (or Docker)
- [ ] Auto-restart enabled
- [ ] Logs configured with rotation
- [ ] Backup system in place
- [ ] Monitoring alerts configured
- [ ] Firewall configured
- [ ] File permissions secured
- [ ] Disk space monitoring enabled
- [ ] Update process documented

---

## Support

For issues:
1. Check logs first
2. Verify yt-dlp works standalone
3. Test with simple commands (`/start`, `/help`)
4. Check GitHub issues
5. Contact support: @support on Telegram
