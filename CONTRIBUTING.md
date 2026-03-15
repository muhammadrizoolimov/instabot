# Contributing to InstaBot2

Thank you for considering contributing to InstaBot2! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community

## How to Contribute

### Reporting Bugs

1. **Check existing issues** to avoid duplicates
2. **Create detailed bug report** with:
   - Description of the issue
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Go version, yt-dlp version)
   - Relevant logs (sanitize sensitive data)

### Suggesting Features

1. **Open a feature request** issue
2. **Describe the feature** clearly:
   - What problem does it solve?
   - How should it work?
   - Any implementation ideas?

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes** following the code style

4. **Test thoroughly**:
   ```bash
   go test ./...
   ```

5. **Commit with clear messages**:
   ```bash
   git commit -m "Add music quality selection feature"
   ```

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request** with:
   - Clear title and description
   - Reference related issues
   - Screenshots/examples if applicable

## Development Setup

### Prerequisites

- Go 1.21+
- yt-dlp
- ffmpeg
- Git

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/instabot.git
cd instabot

# Install dependencies
make setup

# Configure environment
cp .env.example .env
# Add your test bot token

# Run tests
make test

# Run bot
make run
```

## Code Style

### Go Conventions

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting:
  ```bash
  gofmt -w .
  ```
- Use meaningful variable names
- Add comments for exported functions

### Example

```go
// SearchMusic searches for music using yt-dlp with pagination support.
// Returns SearchResponse containing up to 10 results per page.
//
// Parameters:
//   - query: Search query (will be sanitized for UTF-8)
//   - page: Page number (1-based)
//
// Returns error if yt-dlp fails or no results found.
func SearchMusic(query string, page int) (*SearchResponse, error) {
    // Implementation
}
```

### File Organization

```
internal/
├── bot/          # Bot core and handlers
├── music/        # Music-related features
├── services/     # External services (yt-dlp, etc)
├── worker/       # Background job processing
├── database/     # Database operations
├── downloader/   # Platform-specific logic
└── utils/        # Shared utilities
```

### Naming Conventions

- **Files**: lowercase, descriptive (`search.go`, `media.go`)
- **Packages**: lowercase, single word (`music`, `worker`)
- **Types**: PascalCase (`SearchResult`, `WorkerPool`)
- **Functions**: PascalCase for exported, camelCase for private
- **Constants**: UPPER_CASE for constants

## Testing

### Unit Tests

Create tests alongside code:

```go
// music/search_test.go
package music

import "testing"

func TestSearchMusic(t *testing.T) {
    results, err := SearchMusic("test", 1)
    if err != nil {
        t.Fatalf("SearchMusic failed: %v", err)
    }
    
    if len(results.Results) == 0 {
        t.Error("Expected results, got none")
    }
}
```

Run tests:
```bash
go test ./...
```

### Integration Tests

Test with actual Telegram API using test bot:

```go
func TestVideoDownload(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Test with real URL
}
```

## Adding New Features

### Adding a New Platform

1. **Update** `internal/downloader/platforms.go`:
   ```go
   {
       Name: "Facebook",
       Patterns: []string{"facebook.com/", "fb.watch/"},
   }
   ```

2. **Test** platform detection:
   ```go
   func TestDetectFacebook(t *testing.T) {
       platform := DetectPlatform("https://facebook.com/...")
       if platform != "Facebook" {
           t.Errorf("Expected Facebook, got %s", platform)
       }
   }
   ```

3. **Update** README.md with new platform

### Adding a New Command

1. **Update** `internal/bot/handlers.go`:
   ```go
   case strings.HasPrefix(command, "/stats"):
       b.handleStats(message)
   ```

2. **Implement** handler:
   ```go
   func (b *Bot) handleStats(message *tgbotapi.Message) {
       // Implementation
   }
   ```

3. **Update** help text in `/help` command

### Adding Database Tables

1. **Update** schema in `internal/database/db.go`:
   ```go
   CREATE TABLE IF NOT EXISTS user_preferences (
       user_id INTEGER PRIMARY KEY,
       language TEXT DEFAULT 'uz',
       quality TEXT DEFAULT 'best'
   );
   ```

2. **Add** accessor methods:
   ```go
   func (d *Database) SavePreference(userID int64, key, value string) error {
       // Implementation
   }
   ```

## Commit Message Guidelines

### Format

```
<type>: <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting, etc)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

### Examples

```
feat: Add music quality selection

- Users can now choose audio quality (128k, 192k, 320k)
- Quality preference saved in database
- Added quality buttons to music download

Closes #42
```

```
fix: Handle UTF-8 characters in song titles

- Use strings.Map for character filtering
- Properly escape HTML entities
- Truncate long titles to prevent errors

Fixes #38
```

## Documentation

### Code Comments

- Comment all exported functions, types, and constants
- Explain complex logic
- Use examples in godoc comments

### README Updates

- Update feature list when adding features
- Keep installation instructions current
- Update examples if APIs change

### Documentation Files

- Update `docs/ARCHITECTURE.md` for structural changes
- Update `docs/API.md` for API changes
- Update `docs/DEPLOYMENT.md` for deployment changes

## Review Process

1. **Automated checks** must pass:
   - Code compiles
   - Tests pass
   - No linting errors

2. **Maintainer review**:
   - Code quality and style
   - Test coverage
   - Documentation completeness

3. **Feedback addressed**:
   - Respond to review comments
   - Make requested changes
   - Re-request review when ready

4. **Merge**:
   - Squash commits if needed
   - Update changelog
   - Merge to main branch

## Release Process

1. **Version bump** in relevant files
2. **Update** CHANGELOG.md
3. **Tag** release:
   ```bash
   git tag -a v1.2.0 -m "Release v1.2.0"
   git push origin v1.2.0
   ```
4. **Create** GitHub release with notes

## Getting Help

- **Questions**: Open a discussion
- **Bugs**: Create an issue
- **Urgent**: Contact @support on Telegram

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be acknowledged in:
- README.md contributors section
- Release notes
- Project credits

Thank you for contributing to InstaBot2! 🎉
