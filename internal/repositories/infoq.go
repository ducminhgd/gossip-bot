package repositories

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
	"github.com/ducminhgd/gossip-bot/pkg/http"
)

// InfoQRepository handles fetching articles from InfoQ
type InfoQRepository struct {
	httpClient http.HTTPClient
}

// NewInfoQRepository creates a new InfoQRepository
func NewInfoQRepository() *InfoQRepository {
	return &InfoQRepository{
		httpClient: http.NewClient(),
	}
}

// RSS feed structures for InfoQ
type InfoQRSSFeed struct {
	XMLName xml.Name     `xml:"rss"`
	Channel InfoQChannel `xml:"channel"`
}

type InfoQChannel struct {
	Title       string      `xml:"title"`
	Description string      `xml:"description"`
	Link        string      `xml:"link"`
	Items       []InfoQItem `xml:"item"`
}

type InfoQItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Category    string `xml:"category"`
}

// FetchArticles fetches latest articles from InfoQ
func (r *InfoQRepository) FetchArticles(source models.Source) ([]models.News, error) {
	// Fetch RSS feed with appropriate headers for RSS/XML content
	headers := map[string]string{
		"Accept":     "application/rss+xml, application/xml, text/xml, */*",
		"User-Agent": "GossipBot/1.0 RSS Reader (https://github.com/ducminhgd/gossip-bot)",
	}
	body, err := r.httpClient.GetWithHeaders(source.URL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch InfoQ RSS feed: %w", err)
	}

	// Parse RSS feed
	var feed InfoQRSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse InfoQ RSS feed: %w", err)
	}

	// Convert RSS items to news items
	var newsList []models.News
	var skippedArticles []string

	for i, item := range feed.Channel.Items {
		// Limit the number of articles
		if i >= source.Limit {
			break
		}

		// Skip articles with empty titles
		if strings.TrimSpace(item.Title) == "" {
			fmt.Printf("WARNING: skipping InfoQ article with empty title\n")
			skippedArticles = append(skippedArticles, "unknown article")
			continue
		}

		// Parse publication date
		publishedAt, err := r.parseInfoQDate(item.PubDate)
		if err != nil {
			fmt.Printf("WARNING: failed to parse InfoQ article date %s: %v\n", item.PubDate, err)
			publishedAt = time.Now() // Use current time as fallback
		}

		// Clean description (remove HTML tags if any)
		description := r.cleanDescription(item.Description)

		// Create news item
		news := models.News{
			Title:       strings.TrimSpace(item.Title),
			URL:         strings.TrimSpace(item.Link),
			Description: description,
			Source:      "InfoQ",
			SubSource:   strings.TrimSpace(item.Category),
			PublishedAt: publishedAt,
			Score:       0, // InfoQ doesn't provide scores
			Comments:    0, // InfoQ doesn't provide comment counts in RSS
		}

		newsList = append(newsList, news)
	}

	// If all articles were skipped, return an error
	if len(newsList) == 0 && len(skippedArticles) > 0 {
		return nil, fmt.Errorf("failed to fetch any InfoQ articles, skipped: %v", skippedArticles)
	}

	// Sort by publication date (newest first)
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].PublishedAt.After(newsList[j].PublishedAt)
	})

	return newsList, nil
}

// parseInfoQDate parses InfoQ's RSS date format
func (r *InfoQRepository) parseInfoQDate(dateStr string) (time.Time, error) {
	// InfoQ uses RFC1123 format: "Mon, 02 Jan 2006 15:04:05 MST"
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		"Mon, 2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, strings.TrimSpace(dateStr)); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// cleanDescription removes HTML tags and limits description length
func (r *InfoQRepository) cleanDescription(desc string) string {
	// Remove HTML tags (basic cleanup)
	desc = strings.ReplaceAll(desc, "<p>", "")
	desc = strings.ReplaceAll(desc, "</p>", "")
	desc = strings.ReplaceAll(desc, "<br>", " ")
	desc = strings.ReplaceAll(desc, "<br/>", " ")
	desc = strings.ReplaceAll(desc, "<br />", " ")
	desc = strings.ReplaceAll(desc, "&nbsp;", " ")
	desc = strings.ReplaceAll(desc, "&amp;", "&")
	desc = strings.ReplaceAll(desc, "&lt;", "<")
	desc = strings.ReplaceAll(desc, "&gt;", ">")
	desc = strings.ReplaceAll(desc, "&quot;", "\"")
	desc = strings.ReplaceAll(desc, "&#39;", "'")

	// Trim whitespace
	desc = strings.TrimSpace(desc)

	// Limit description length
	if len(desc) > 200 {
		desc = desc[:200] + "..."
	}

	return desc
}
