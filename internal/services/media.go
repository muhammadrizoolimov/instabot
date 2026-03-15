package services

import (
	"fmt"
	"instabot2/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type MediaService struct {
	TempDir string
}

type DownloadResult struct {
	FilePath  string
	Title     string
	Thumbnail string
	Duration  int
}

func NewMediaService(tempDir string) *MediaService {
	if tempDir == "" {
		tempDir = "./temp"
	}
	os.MkdirAll(tempDir, 0755)
	return &MediaService{TempDir: tempDir}
}

// DownloadVideo downloads video from URL using yt-dlp
func (ms *MediaService) DownloadVideo(url string) (*DownloadResult, error) {
	timestamp := time.Now().UnixNano()
	outputTemplate := filepath.Join(ms.TempDir, fmt.Sprintf("video_%d.%%(ext)s", timestamp))

	cmd := exec.Command("yt-dlp",
		"-f", "best[ext=mp4]/best",
		"--no-playlist",
		"-o", outputTemplate,
		"--print", "after_move:filepath",
		"--print", "title",
		"--print", "duration",
		url,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("download failed: %v, output: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected yt-dlp output")
	}

	filePath := strings.TrimSpace(lines[len(lines)-3])
	title := "Downloaded Video"
	duration := 0

	if len(lines) >= 2 {
		title = utils.SanitizeForTelegram(strings.TrimSpace(lines[len(lines)-2]))
	}
	if len(lines) >= 1 {
		fmt.Sscanf(strings.TrimSpace(lines[len(lines)-1]), "%d", &duration)
	}

	return &DownloadResult{
		FilePath: filePath,
		Title:    title,
		Duration: duration,
	}, nil
}

// ExtractAudio extracts audio from video file or downloads audio directly
func (ms *MediaService) ExtractAudio(url string) (*DownloadResult, error) {
	timestamp := time.Now().UnixNano()
	outputTemplate := filepath.Join(ms.TempDir, fmt.Sprintf("audio_%d.%%(ext)s", timestamp))

	cmd := exec.Command("yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--no-playlist",
		"-o", outputTemplate,
		"--print", "after_move:filepath",
		"--print", "title",
		"--print", "duration",
		"--print", "uploader",
		url,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("audio extraction failed: %v, output: %s", err, string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("unexpected yt-dlp output")
	}

	filePath := ""
	title := "Downloaded Audio"
	duration := 0

	// Parse output lines
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ms.TempDir) && strings.HasSuffix(line, ".mp3") {
			filePath = line
		} else if i == len(lines)-3 {
			title = utils.SanitizeForTelegram(line)
		} else if i == len(lines)-2 {
			fmt.Sscanf(line, "%d", &duration)
		}
	}

	if filePath == "" {
		return nil, fmt.Errorf("audio file not found in output")
	}

	return &DownloadResult{
		FilePath: filePath,
		Title:    title,
		Duration: duration,
	}, nil
}

// Cleanup removes temporary file
func (ms *MediaService) Cleanup(filePath string) error {
	if filePath != "" && strings.Contains(filePath, ms.TempDir) {
		return os.Remove(filePath)
	}
	return nil
}
