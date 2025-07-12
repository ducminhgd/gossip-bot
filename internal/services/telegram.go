package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	TELEGRAM_SEND_MESSAGE_URL_PATTERN = "https://api.telegram.org/bot%s/sendMessage"

	TELEGRAM_PARSE_MODE_DEFAULT  = "Markdown"
	TELEGRAM_PARSE_MODE_MARKDOWN = "Markdown"
)

type TelegramService struct {
	httpClient *http.Client
	botToken   string
}

func NewTelegramService(botToken string) *TelegramService {
	return &TelegramService{
		botToken:   botToken,
		httpClient: http.DefaultClient,
	}
}

func (s *TelegramService) SendMessage(message string, chat_id int64, thread_id int64, parse_mode string) error {
	url := fmt.Sprintf(TELEGRAM_SEND_MESSAGE_URL_PATTERN, s.botToken)

	body := map[string]interface{}{
		"chat_id":           chat_id,
		"text":              message,
		"message_thread_id": thread_id,
		"parse_mode":        parse_mode,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Failed to marshal body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
