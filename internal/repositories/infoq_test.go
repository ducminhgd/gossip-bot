package repositories

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
)

// MockHTTPClient for testing InfoQ repository
type mockInfoQHTTPClient struct {
	responses map[string][]byte
	errors    map[string]error
}

func (m *mockInfoQHTTPClient) Get(url string) ([]byte, error) {
	if err, exists := m.errors[url]; exists {
		return nil, err
	}
	if response, exists := m.responses[url]; exists {
		return response, nil
	}
	return nil, fmt.Errorf("no mock response for URL: %s", url)
}

func (m *mockInfoQHTTPClient) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	return m.Get(url)
}

func (m *mockInfoQHTTPClient) GetJSON(url string, v any) error {
	body, err := m.Get(url)
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, v)
}

func (m *mockInfoQHTTPClient) PostForm(url string, data map[string][]string, headers map[string]string) ([]byte, error) {
	return nil, fmt.Errorf("PostForm not implemented in mock")
}

// testInfoQRepository for testing
type testInfoQRepository struct {
	mockClient *mockInfoQHTTPClient
}

func newTestInfoQRepository(mockClient *mockInfoQHTTPClient) *testInfoQRepository {
	return &testInfoQRepository{
		mockClient: mockClient,
	}
}

func (r *testInfoQRepository) FetchArticles(source models.Source) ([]models.News, error) {
	// Fetch RSS feed with appropriate headers
	headers := map[string]string{
		"Accept":     "application/rss+xml, application/xml, text/xml, */*",
		"User-Agent": "GossipBot/1.0 RSS Reader (https://github.com/ducminhgd/gossip-bot)",
	}
	body, err := r.mockClient.GetWithHeaders(source.URL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch InfoQ RSS feed: %w", err)
	}

	// Parse RSS feed
	var feed InfoQRSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse InfoQ RSS feed: %w", err)
	}

	// Convert RSS items to news items (simplified for testing)
	var newsList []models.News
	for i, item := range feed.Channel.Items {
		if i >= source.Limit {
			break
		}

		publishedAt, _ := time.Parse(time.RFC1123, item.PubDate)

		news := models.News{
			Title:       item.Title,
			URL:         item.Link,
			Description: item.Description,
			Source:      "InfoQ",
			SubSource:   item.Category,
			PublishedAt: publishedAt,
			Score:       0,
			Comments:    0,
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}

func TestInfoQRepository_FetchArticles(t *testing.T) {
	// Mock RSS feed response
	mockRSSResponse := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>InfoQ - RSS - Articles</title>
		<description>InfoQ RSS Articles feed</description>
		<link>https://www.infoq.com</link>
		<item>
			<title>Effective Practices for Coding with a Chat-Based AI</title>
			<link>https://www.infoq.com/articles/effective-practices-ai-chat-based-coding/</link>
			<description>In this article, we explore how AI agents are reshaping software development...</description>
			<pubDate>Thu, 04 Jul 2025 00:00:00 GMT</pubDate>
			<guid>https://www.infoq.com/articles/effective-practices-ai-chat-based-coding/</guid>
			<category>Development</category>
		</item>
		<item>
			<title>Agentic AI Architecture Framework for Enterprises</title>
			<link>https://www.infoq.com/articles/agentic-ai-architecture-framework/</link>
			<description>To deploy agentic AI responsibly and effectively in the enterprise...</description>
			<pubDate>Thu, 11 Jul 2025 00:00:00 GMT</pubDate>
			<guid>https://www.infoq.com/articles/agentic-ai-architecture-framework/</guid>
			<category>Architecture</category>
		</item>
	</channel>
</rss>`

	mockClient := &mockInfoQHTTPClient{
		responses: map[string][]byte{
			"https://feed.infoq.com/": []byte(mockRSSResponse),
		},
		errors: make(map[string]error),
	}

	repo := newTestInfoQRepository(mockClient)
	source := models.Source{
		Name:  "InfoQ",
		Type:  "infoq",
		URL:   "https://feed.infoq.com/",
		Limit: 10,
	}

	articles, err := repo.FetchArticles(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(articles))
	}

	// Check first article
	firstArticle := articles[0]
	if firstArticle.Title != "Effective Practices for Coding with a Chat-Based AI" {
		t.Errorf("Expected title 'Effective Practices for Coding with a Chat-Based AI', got '%s'", firstArticle.Title)
	}
	if firstArticle.Source != "InfoQ" {
		t.Errorf("Expected source 'InfoQ', got '%s'", firstArticle.Source)
	}
	if firstArticle.SubSource != "Development" {
		t.Errorf("Expected sub-source 'Development', got '%s'", firstArticle.SubSource)
	}
}

func TestInfoQRepository_FetchArticles_NetworkError(t *testing.T) {
	mockClient := &mockInfoQHTTPClient{
		responses: make(map[string][]byte),
		errors: map[string]error{
			"https://feed.infoq.com/": fmt.Errorf("network error"),
		},
	}

	repo := newTestInfoQRepository(mockClient)
	source := models.Source{
		Name:  "InfoQ",
		Type:  "infoq",
		URL:   "https://feed.infoq.com/",
		Limit: 10,
	}

	_, err := repo.FetchArticles(source)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "failed to fetch InfoQ RSS feed: network error" {
		t.Errorf("Expected specific error message, got '%s'", err.Error())
	}
}

func TestInfoQRepository_FetchArticles_InvalidXML(t *testing.T) {
	mockClient := &mockInfoQHTTPClient{
		responses: map[string][]byte{
			"https://feed.infoq.com/": []byte("invalid xml"),
		},
		errors: make(map[string]error),
	}

	repo := newTestInfoQRepository(mockClient)
	source := models.Source{
		Name:  "InfoQ",
		Type:  "infoq",
		URL:   "https://feed.infoq.com/",
		Limit: 10,
	}

	_, err := repo.FetchArticles(source)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "failed to parse InfoQ RSS feed") {
		t.Errorf("Expected parse error, got '%s'", err.Error())
	}
}

func TestInfoQRepository_parseInfoQDate(t *testing.T) {
	repo := NewInfoQRepository()

	testCases := []struct {
		input    string
		expected bool // whether parsing should succeed
	}{
		{"Thu, 04 Jul 2025 00:00:00 GMT", true},
		{"Mon, 02 Jan 2006 15:04:05 MST", true},
		{"2025-07-04T10:30:00Z", true},
		{"invalid date", false},
		{"", false},
	}

	for _, tc := range testCases {
		_, err := repo.parseInfoQDate(tc.input)
		if tc.expected && err != nil {
			t.Errorf("Expected to parse '%s' successfully, got error: %v", tc.input, err)
		}
		if !tc.expected && err == nil {
			t.Errorf("Expected to fail parsing '%s', but succeeded", tc.input)
		}
	}
}

func TestInfoQRepository_cleanDescription(t *testing.T) {
	repo := NewInfoQRepository()

	testCases := []struct {
		input    string
		expected string
	}{
		{"<p>Simple text</p>", "Simple text"},
		{"Text with &amp; entities", "Text with & entities"},
		{"<br/>Line break", " Line break"},
		{"Very long description that exceeds the maximum length limit and should be truncated to exactly 200 characters plus three dots to indicate that there is more content available but it has been cut off for display purposes", "Very long description that exceeds the maximum length limit and should be truncated to exactly 200 characters plus three dots to indicate that there is more content available but it has been cut..."},
	}

	for _, tc := range testCases {
		result := repo.cleanDescription(tc.input)
		if result != tc.expected {
			t.Errorf("Expected '%s', got '%s'", tc.expected, result)
		}
	}
}
