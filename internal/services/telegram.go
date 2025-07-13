package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GossipBot/1.0 (https://github.com/ducminhgd/gossip-bot)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		b, _ := io.ReadAll(resp.Body)
		fmt.Printf("send Telegram failed: %s\n", string(b))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
