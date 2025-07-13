package main

import (
	"fmt"
	"log"
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

	// Generate issue content
	log.Println("Generating issue content...")
	issueContent, err := githubService.GenerateIssueContent(newsMap)
	if err != nil {
		log.Fatalf("Failed to generate issue content: %v", err)
	}

	// Create issue
	now := time.Now().UTC()
	issueTitle := fmt.Sprintf("Daily News Digest - %s", now.Format("2006-01-02"))

	log.Printf("Creating GitHub issue: %s", issueTitle)
	issue, err := githubService.CreateIssue(issueTitle, issueContent)
	if err != nil {
		log.Fatalf("Failed to create issue: %v", err)
	}

	log.Printf("Successfully created issue #%d: %s", issue.GetNumber(), issue.GetHTMLURL())

	// Send Telegram message
	telegramCfg, err := config.LoadTelegramConfig()
	if err != nil {
		log.Fatalf("Failed to load Telegram configuration: %v", err)
	}

	telegramService := services.NewTelegramService(telegramCfg.TelegramBotToken)
	err = telegramService.SendMessage(issueContent, telegramCfg.TelegramChatID, telegramCfg.TelegramThreadID, services.TELEGRAM_PARSE_MODE_MARKDOWN)
	if err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
	} else {
		log.Println("Successfully sent Telegram message")
	}
}
