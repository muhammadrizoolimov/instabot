package music

import (
	"encoding/json"
	"fmt"
	"instabot2/internal/utils"
	"os/exec"
	"strings"
)

const (
	ResultsPerPage = 10
)

type SearchResult struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Duration int    `json:"duration"`
	Uploader string `json:"uploader"`
}

type SearchResponse struct {
	Results    []SearchResult
	Page       int
	TotalPages int
}

// SearchMusic searches for music using yt-dlp
func SearchMusic(query string, page int) (*SearchResponse, error) {
	if page < 1 {
		page = 1
	}

	// Clean query for UTF-8 safety
	query = utils.SanitizeForTelegram(query)
	if query == "" {
		return nil, fmt.Errorf("empty search query")
	}

	// Use yt-dlp to search on YouTube Music
	searchQuery := fmt.Sprintf("ytsearch%d:%s", ResultsPerPage, query)

	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--no-playlist",
		"--flat-playlist",
		"--skip-download",
		searchQuery,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp search failed: %v, output: %s", err, string(output))
	}

	// Parse results (yt-dlp outputs one JSON object per line)
	lines := strings.Split(string(output), "\n")
	var results []SearchResult

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		// Extract fields with safe type assertions
		title := ""
		if t, ok := result["title"].(string); ok {
			title = utils.SanitizeForTelegram(t)
		}

		url := ""
		if u, ok := result["url"].(string); ok {
			url = u
		} else if u, ok := result["webpage_url"].(string); ok {
			url = u
		} else if id, ok := result["id"].(string); ok {
			url = fmt.Sprintf("https://www.youtube.com/watch?v=%s", id)
		}

		duration := 0
		if d, ok := result["duration"].(float64); ok {
			duration = int(d)
		}

		uploader := ""
		if up, ok := result["uploader"].(string); ok {
			uploader = utils.SanitizeForTelegram(up)
		}

		if title != "" && url != "" {
			results = append(results, SearchResult{
				Title:    title,
				URL:      url,
				Duration: duration,
				Uploader: uploader,
			})
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for: %s", query)
	}

	// Calculate pagination
	totalResults := len(results)
	startIdx := (page - 1) * ResultsPerPage
	endIdx := startIdx + ResultsPerPage

	if startIdx >= totalResults {
		startIdx = 0
		endIdx = ResultsPerPage
		page = 1
	}

	if endIdx > totalResults {
		endIdx = totalResults
	}

	pageResults := results[startIdx:endIdx]
	totalPages := (totalResults + ResultsPerPage - 1) / ResultsPerPage

	return &SearchResponse{
		Results:    pageResults,
		Page:       page,
		TotalPages: totalPages,
	}, nil
}

// FormatDuration converts seconds to MM:SS format
func FormatDuration(seconds int) string {
	if seconds == 0 {
		return "Unknown"
	}
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, secs)
}
