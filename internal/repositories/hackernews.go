package repositories

import (
	"fmt"
	"sort"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
	"github.com/ducminhgd/gossip-bot/pkg/http"
)

// HackerNewsRepository handles fetching news from Hacker News
type HackerNewsRepository struct {
	httpClient *http.Client
}

// NewHackerNewsRepository creates a new HackerNewsRepository
func NewHackerNewsRepository() *HackerNewsRepository {
	return &HackerNewsRepository{
		httpClient: http.NewClient(),
	}
}

// FetchTopStories fetches top stories from Hacker News
func (r *HackerNewsRepository) FetchTopStories(source models.Source) ([]models.News, error) {
	// Hacker News API uses Firebase API
	// The base URL should be https://hacker-news.firebaseio.com/v0
	// But the actual website is https://news.ycombinator.com/

	// Fetch top stories IDs
	topStoriesURL := "https://hacker-news.firebaseio.com/v0/topstories.json"
	var storyIDs []int
	if err := r.httpClient.GetJSON(topStoriesURL, &storyIDs); err != nil {
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
		story, err := r.fetchStory(id)
		if err != nil {
			// Log warning and continue with other stories
			fmt.Printf("WARNING: failed to fetch Hacker News story %d: %v\n", id, err)
			skippedStories = append(skippedStories, id)
			continue
		}
		newsList = append(newsList, story)
	}

	// If all stories were skipped, return an error
	if len(newsList) == 0 && len(skippedStories) > 0 {
		return nil, fmt.Errorf("failed to fetch any Hacker News top stories, skipped: %v", skippedStories)
	}

	// Sort by score
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].Score > newsList[j].Score
	})

	return newsList, nil
}

// FetchBestStories fetches best stories from Hacker News
func (r *HackerNewsRepository) FetchBestStories(source models.Source) ([]models.News, error) {
	// Fetch best stories IDs
	bestStoriesURL := "https://hacker-news.firebaseio.com/v0/beststories.json"
	var storyIDs []int
	if err := r.httpClient.GetJSON(bestStoriesURL, &storyIDs); err != nil {
		return nil, fmt.Errorf("failed to fetch best stories: %w", err)
	}

	// Limit the number of stories
	if len(storyIDs) > source.Limit {
		storyIDs = storyIDs[:source.Limit]
	}

	// Fetch each story
	var newsList []models.News
	var skippedStories []int

	for _, id := range storyIDs {
		story, err := r.fetchStory(id)
		if err != nil {
			// Log warning and continue with other stories
			fmt.Printf("WARNING: failed to fetch Hacker News story %d: %v\n", id, err)
			skippedStories = append(skippedStories, id)
			continue
		}
		newsList = append(newsList, story)
	}

	// If all stories were skipped, return an error
	if len(newsList) == 0 && len(skippedStories) > 0 {
		return nil, fmt.Errorf("failed to fetch any Hacker News best stories, skipped: %v", skippedStories)
	}

	// Sort by score
	sort.Slice(newsList, func(i, j int) bool {
		return newsList[i].Score > newsList[j].Score
	})

	return newsList, nil
}

// fetchStory fetches a single story from Hacker News
func (r *HackerNewsRepository) fetchStory(id int) (models.News, error) {
	storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	var story map[string]any
	if err := r.httpClient.GetJSON(storyURL, &story); err != nil {
		return models.News{}, fmt.Errorf("failed to fetch story %d: %w", id, err)
	}

	// Extract story details
	title, _ := story["title"].(string)
	url, _ := story["url"].(string)
	if url == "" {
		// If the URL is empty, it's a self-post, so use the Hacker News item URL
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

	return news, nil
}
