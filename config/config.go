package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ducminhgd/gossip-bot/internal/models"
	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	// GithubToken is the GitHub token used to create issues
	GithubToken string

	// GithubOwner is the owner of the GitHub repository
	GithubOwner string

	// GithubRepo is the name of the GitHub repository
	GithubRepo string

	// Sources is a list of news sources
	Sources []models.Source
}

type TelegramConfig struct {
	// TelegramBotToken is the Telegram bot token used to send messages
	TelegramBotToken string

	// TelegramChatID is the chat ID to send messages to
	TelegramChatID int64

	// TelegramThreadID is the thread ID to send messages to
	TelegramThreadID int64
}

type RedditAppConfig struct {
	AppID     string
	AppSecret string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get GitHub configuration
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	githubOwner := os.Getenv("GITHUB_OWNER")
	if githubOwner == "" {
		return nil, fmt.Errorf("GITHUB_OWNER environment variable is required")
	}

	githubRepo := os.Getenv("GITHUB_REPO")
	if githubRepo == "" {
		return nil, fmt.Errorf("GITHUB_REPO environment variable is required")
	}

	// Get sources configuration
	sourcesList := os.Getenv("SOURCES")
	if sourcesList == "" {
		return nil, fmt.Errorf("SOURCES environment variable is required")
	}

	sources := []models.Source{}
	sourceNames := strings.Split(sourcesList, ",")

	for _, sourceName := range sourceNames {
		sourceName = strings.TrimSpace(sourceName)
		if sourceName == "" {
			continue
		}

		sourceType := os.Getenv(fmt.Sprintf("SOURCE_%s_TYPE", sourceName))
		if sourceType == "" {
			return nil, fmt.Errorf("SOURCE_%s_TYPE environment variable is required", sourceName)
		}

		sourceURL := os.Getenv(fmt.Sprintf("SOURCE_%s_URL", sourceName))
		if sourceURL == "" {
			return nil, fmt.Errorf("SOURCE_%s_URL environment variable is required", sourceName)
		}

		sourceLimitStr := os.Getenv(fmt.Sprintf("SOURCE_%s_LIMIT", sourceName))
		if sourceLimitStr == "" {
			sourceLimitStr = "10" // Default limit
		}

		sourceLimit, err := strconv.Atoi(sourceLimitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SOURCE_%s_LIMIT: %v", sourceName, err)
		}

		sourceSubSource := os.Getenv(fmt.Sprintf("SOURCE_%s_SUBSOURCE", sourceName))

		source := models.Source{
			Name:      sourceName,
			Type:      sourceType,
			URL:       sourceURL,
			Limit:     sourceLimit,
			SubSource: sourceSubSource,
		}

		sources = append(sources, source)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no valid sources configured")
	}

	return &Config{
		GithubToken: githubToken,
		GithubOwner: githubOwner,
		GithubRepo:  githubRepo,
		Sources:     sources,
	}, nil
}

func LoadTelegramConfig() (*TelegramConfig, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get Telegram configuration
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	telegramChatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramChatIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_CHAT_ID environment variable is required")
	}

	telegramChatID, err := strconv.ParseInt(telegramChatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TELEGRAM_CHAT_ID: %v", err)
	}

	telegramThreadIDStr := os.Getenv("TELEGRAM_THREAD_ID")
	if telegramThreadIDStr == "" {
		return nil, fmt.Errorf("TELEGRAM_THREAD_ID environment variable is required")
	}

	telegramThreadID, err := strconv.ParseInt(telegramThreadIDStr, 10, 64)
	if err != nil {
		telegramThreadID = 0
	}

	return &TelegramConfig{
		TelegramBotToken: telegramBotToken,
		TelegramChatID:   telegramChatID,
		TelegramThreadID: telegramThreadID,
	}, nil
}

func LoadRedditAppConfig() (*RedditAppConfig, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get Reddit app configuration
	redditAppID := os.Getenv("REDDIT_APP_ID")
	if redditAppID == "" {
		return nil, fmt.Errorf("REDDIT_APP_ID environment variable is required")
	}

	redditAppSecret := os.Getenv("REDDIT_APP_SECRET")
	if redditAppSecret == "" {
		return nil, fmt.Errorf("REDDIT_APP_SECRET environment variable is required")
	}

	return &RedditAppConfig{
		AppID:     redditAppID,
		AppSecret: redditAppSecret,
	}, nil
}
