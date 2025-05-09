package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
	"github.com/ducminhgd/gossip-bot/pkg/http"
)

// NewsService handles operations related to news
type NewsService struct {
	httpClient *http.Client
	sources    []models.Source
}

// NewNewsService creates a new NewsService
func NewNewsService(sources []models.Source) *NewsService {
	return &NewsService{
		httpClient: http.NewClient(),
		sources:    sources,
	}
}

// FetchAllNews fetches news from all sources
// If a source can't be crawled, it will be skipped and a warning will be logged
func (s *NewsService) FetchAllNews() (map[string][]models.News, error) {
	result := make(map[string][]models.News)
	var skippedSources []string

	for _, source := range s.sources {
		news, err := s.FetchNewsBySource(source)
		if err != nil {
			// Log warning and continue with other sources
			fmt.Printf("WARNING: failed to fetch news from %s: %v\n", source.Name, err)
			skippedSources = append(skippedSources, source.Name)
			continue
		}
		result[source.Name] = news
	}

	// If all sources were skipped, return an error
	if len(result) == 0 && len(skippedSources) > 0 {
		return nil, fmt.Errorf("failed to fetch news from any source, skipped: %v", skippedSources)
	}

	return result, nil
}

// FetchNewsBySource fetches news from a specific source
func (s *NewsService) FetchNewsBySource(source models.Source) ([]models.News, error) {
	switch strings.ToLower(source.Type) {
	case "hackernews":
		return s.fetchHackerNews(source)
	case "reddit":
		return s.fetchReddit(source)
	default:
		return nil, fmt.Errorf("unsupported source type: %s", source.Type)
	}
}

// fetchHackerNews fetches news from Hacker News
func (s *NewsService) fetchHackerNews(source models.Source) ([]models.News, error) {
	// Fetch top stories
	topStoriesURL := "https://hacker-news.firebaseio.com/v0/topstories.json"
	var storyIDs []int
	if err := s.httpClient.GetJSON(topStoriesURL, &storyIDs); err != nil {
		return nil, fmt.Errorf("failed to fetch top stories: %w", err)
	}

	// Limit the number of stories
	if len(storyIDs) > source.Limit {
		storyIDs = storyIDs[:source.Limit]
	}

	// Fetch each story
	var newsList []models.News
	var skippedStories []int

	for _, id := range storyIDs {
		storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
		var story map[string]any
		if err := s.httpClient.GetJSON(storyURL, &story); err != nil {
			// Log warning and continue with other stories
			fmt.Printf("WARNING: failed to fetch Hacker News story %d: %v\n", id, err)
			skippedStories = append(skippedStories, id)
			continue
		}

		// Extract story details
		title, _ := story["title"].(string)
		url, _ := story["url"].(string)
		if url == "" {
			url = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", id)
		}
		score, _ := story["score"].(float64)
		descendants, _ := story["descendants"].(float64)
		unixTime, _ := story["time"].(float64)
		publishedAt := time.Unix(int64(unixTime), 0)

		// Create news item
		news := models.News{
			Title:       title,
			URL:         url,
			Description: fmt.Sprintf("Score: %d, Comments: %d", int(score), int(descendants)),
			Source:      "Hacker News",
			PublishedAt: publishedAt,
			Score:       int(score),
			Comments:    int(descendants),
		}

		newsList = append(newsList, news)
	}

	// If all stories were skipped, return an error
	if len(newsList) == 0 && len(skippedStories) > 0 {
		return nil, fmt.Errorf("failed to fetch any Hacker News stories, skipped: %v", skippedStories)
	}

	// Sort by score
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].Score > newsList[j].Score
	})

	return newsList, nil
}

// fetchReddit fetches news from Reddit
func (s *NewsService) fetchReddit(source models.Source) ([]models.News, error) {
	// Construct URL
	subreddit := source.SubSource
	if subreddit == "" {
		return nil, fmt.Errorf("subreddit is required for Reddit source")
	}

	redditURL := fmt.Sprintf("%s/r/%s/hot.json?limit=%d", source.URL, url.PathEscape(subreddit), source.Limit)

	// Fetch data
	body, err := s.httpClient.Get(redditURL)
	if err != nil {
		// If we get a 403 Forbidden error, it's likely due to Reddit's API restrictions
		// This is common in CI/CD environments like GitHub Actions
		if strings.Contains(err.Error(), "403") {
			return nil, fmt.Errorf("Reddit API access forbidden (403) - this is common in CI/CD environments: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch Reddit data: %w", err)
	}

	// Parse response
	var response struct {
		Data struct {
			Children []struct {
				Data struct {
					Title       string  `json:"title"`
					URL         string  `json:"url"`
					Permalink   string  `json:"permalink"`
					Score       int     `json:"score"`
					NumComments int     `json:"num_comments"`
					Created     float64 `json:"created_utc"`
					Selftext    string  `json:"selftext"`
				} `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse Reddit response: %w", err)
	}

	// Extract posts
	var newsList []models.News
	var skippedPosts []string

	for _, child := range response.Data.Children {
		post := child.Data

		// Skip stickied posts or announcements
		if strings.HasPrefix(strings.ToLower(post.Title), "[announcement]") {
			continue
		}

		// Skip posts with empty titles (shouldn't happen, but just in case)
		if post.Title == "" {
			fmt.Printf("WARNING: skipping Reddit post with empty title in r/%s\n", subreddit)
			skippedPosts = append(skippedPosts, "unknown post")
			continue
		}

		// Create URL (use permalink if URL is empty)
		postURL := post.URL
		if postURL == "" || strings.HasPrefix(postURL, "/r/") {
			postURL = fmt.Sprintf("https://www.reddit.com%s", post.Permalink)
		}

		// Create description
		description := post.Selftext
		if len(description) > 100 {
			description = description[:100] + "..."
		}
		if description == "" {
			description = fmt.Sprintf("Score: %d, Comments: %d", post.Score, post.NumComments)
		}

		// Create news item
		news := models.News{
			Title:       post.Title,
			URL:         postURL,
			Description: description,
			Source:      "Reddit",
			SubSource:   subreddit,
			PublishedAt: time.Unix(int64(post.Created), 0),
			Score:       post.Score,
			Comments:    post.NumComments,
		}

		newsList = append(newsList, news)
	}

	// If all posts were skipped, return an error
	if len(newsList) == 0 && len(skippedPosts) > 0 {
		return nil, fmt.Errorf("failed to fetch any Reddit posts from r/%s, skipped: %v", subreddit, skippedPosts)
	}

	// Sort by score
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].Score > newsList[j].Score
	})

	return newsList, nil
}
