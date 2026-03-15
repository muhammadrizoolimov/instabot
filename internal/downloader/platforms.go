package downloader

import (
	"strings"
)

// Platform represents a supported social media platform
type Platform struct {
	Name     string
	Patterns []string
}

var Platforms = []Platform{
	{
		Name: "Instagram",
		Patterns: []string{
			"instagram.com/",
			"instagr.am/",
		},
	},
	{
		Name: "TikTok",
		Patterns: []string{
			"tiktok.com/",
			"vm.tiktok.com/",
			"vt.tiktok.com/",
		},
	},
	{
		Name: "YouTube",
		Patterns: []string{
			"youtube.com/",
			"youtu.be/",
			"youtube.com/shorts/",
		},
	},
	{
		Name: "Pinterest",
		Patterns: []string{
			"pinterest.com/",
			"pin.it/",
		},
	},
	{
		Name: "Snapchat",
		Patterns: []string{
			"snapchat.com/",
		},
	},
	{
		Name: "Likee",
		Patterns: []string{
			"likee.video/",
			"like.video/",
		},
	},
}

// DetectPlatform detects which platform a URL belongs to
func DetectPlatform(url string) string {
	url = strings.ToLower(url)
	for _, platform := range Platforms {
		for _, pattern := range platform.Patterns {
			if strings.Contains(url, pattern) {
				return platform.Name
			}
		}
	}
	return "Unknown"
}

// IsSupportedURL checks if URL is from a supported platform
func IsSupportedURL(url string) bool {
	return DetectPlatform(url) != "Unknown"
}
