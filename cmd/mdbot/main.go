package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ducminhgd/gossip-bot/config"
	"github.com/ducminhgd/gossip-bot/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create services
	newsService := services.NewNewsService(cfg.Sources)
	githubService := services.NewGithubService(
		cfg.GithubToken,
		cfg.GithubOwner,
		cfg.GithubRepo,
	)

	// Fetch news from all sources
	log.Println("Fetching news from all sources...")
	newsMap, err := newsService.FetchAllNews()
	if err != nil {
		log.Fatalf("Failed to fetch news: %v", err)
	}

	// Generate markdown content
	log.Println("Generating markdown content...")
	markdownContent, err := githubService.GenerateIssueContent(newsMap)
	if err != nil {
		log.Fatalf("Failed to generate markdown content: %v", err)
	}

	// Create news directory if it doesn't exist
	newsDir := "news"
	if _, err := os.Stat(newsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(newsDir, 0755); err != nil {
			log.Fatalf("Failed to create news directory: %v", err)
		}
	}

	// Create markdown file
	now := time.Now().UTC()
	filename := now.Format("2006-01-02") + ".md"
	filePath := filepath.Join(newsDir, filename)

	log.Printf("Creating markdown file: %s", filePath)
	if err := os.WriteFile(filePath, []byte(markdownContent), 0644); err != nil {
		log.Fatalf("Failed to write markdown file: %v", err)
	}

	log.Printf("Successfully created markdown file: %s", filePath)

	// Send Telegram message
	// telegramCfg, err := config.LoadTelegramConfig()
	// if err != nil {
	// 	log.Fatalf("Failed to load Telegram configuration: %v", err)
	// }

	// telegramService := services.NewTelegramService(telegramCfg.TelegramBotToken)
	// err = telegramService.SendMessage(markdownContent, telegramCfg.TelegramChatID, telegramCfg.TelegramThreadID, services.TELEGRAM_PARSE_MODE_MARKDOWN)
	// if err != nil {
	// 	log.Fatalf("Failed to send Telegram message: %v", err)
	// }
	// log.Println("Successfully sent Telegram message")
}
