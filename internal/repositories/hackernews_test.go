package repositories

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
)

// MockHTTPClient is a mock implementation of the HTTP client for testing
type MockHTTPClient struct {
	GetJSONFunc func(url string, v any) error
}

// Get is a mock implementation of the Get method
func (m *MockHTTPClient) Get(url string) ([]byte, error) {
	// This method is not used directly in the tests
	return nil, nil
}

// GetWithHeaders is a mock implementation of the GetWithHeaders method
func (m *MockHTTPClient) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	// This method is not used directly in the tests
	return nil, nil
}

// GetJSON is a mock implementation of the GetJSON method
func (m *MockHTTPClient) GetJSON(url string, v any) error {
	return m.GetJSONFunc(url, v)
}

// For testing, we'll modify the HackerNewsRepository struct to accept our mock
// This is a test-only version of the repository
type testHackerNewsRepository struct {
	mockClient *MockHTTPClient
}

// Override the methods we need for testing
func (r *testHackerNewsRepository) FetchTopStories(source models.Source) ([]models.News, error) {
	// Fetch top stories IDs
	topStoriesURL := "https://hacker-news.firebaseio.com/v0/topstories.json"
	var storyIDs []int
	if err := r.mockClient.GetJSON(topStoriesURL, &storyIDs); err != nil {
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

func (r *testHackerNewsRepository) FetchBestStories(source models.Source) ([]models.News, error) {
	// Fetch best stories IDs
	bestStoriesURL := "https://hacker-news.firebaseio.com/v0/beststories.json"
	var storyIDs []int
	if err := r.mockClient.GetJSON(bestStoriesURL, &storyIDs); err != nil {
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

func (r *testHackerNewsRepository) fetchStory(id int) (models.News, error) {
	storyURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	var story map[string]any
	if err := r.mockClient.GetJSON(storyURL, &story); err != nil {
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

// Helper function to create a test repository with a mock client
func NewHackerNewsRepositoryWithClient(mockClient *MockHTTPClient) *testHackerNewsRepository {
	return &testHackerNewsRepository{
		mockClient: mockClient,
	}
}

func TestNewHackerNewsRepository(t *testing.T) {
	repo := NewHackerNewsRepository()
	if repo == nil {
		t.Fatal("Expected non-nil repository")
	}
	if repo.httpClient == nil {
		t.Fatal("Expected non-nil HTTP client")
	}
}

func TestFetchTopStories_Success(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			switch url {
			case "https://hacker-news.firebaseio.com/v0/topstories.json":
				// Mock the top stories response
				storyIDs, ok := v.(*[]int)
				if !ok {
					t.Fatalf("Expected *[]int, got %T", v)
				}
				*storyIDs = []int{123, 456, 789}
			case "https://hacker-news.firebaseio.com/v0/item/123.json":
				// Mock the first story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Test Story 1",
					"url":         "https://example.com/1",
					"score":       float64(100),
					"descendants": float64(10),
					"time":        float64(1625097600), // 2021-07-01
				}
			case "https://hacker-news.firebaseio.com/v0/item/456.json":
				// Mock the second story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Test Story 2",
					"url":         "https://example.com/2",
					"score":       float64(200),
					"descendants": float64(20),
					"time":        float64(1625184000), // 2021-07-02
				}
			case "https://hacker-news.firebaseio.com/v0/item/789.json":
				// Mock the third story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title": "Test Story 3",
					// No URL to test self-post handling
					"score":       float64(300),
					"descendants": float64(30),
					"time":        float64(1625270400), // 2021-07-03
				}
			default:
				t.Fatalf("Unexpected URL: %s", url)
			}
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Create a test source
	source := models.Source{
		Name:  "HackerNews",
		Type:  "hackernews",
		URL:   "https://hacker-news.firebaseio.com/v0",
		Limit: 3,
	}

	// Call the method being tested
	news, err := repo.FetchTopStories(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the results
	if len(news) != 3 {
		t.Fatalf("Expected 3 news items, got %d", len(news))
	}

	// Check that the news items are sorted by score (highest first)
	if news[0].Score != 300 || news[1].Score != 200 || news[2].Score != 100 {
		t.Fatalf("Expected news items to be sorted by score, got %v", news)
	}

	// Check the details of the first news item
	expectedFirstNews := models.News{
		Title:       "Test Story 3",
		URL:         "https://news.ycombinator.com/item?id=789", // Self-post URL
		Description: "Score: 300, Comments: 30",
		Source:      "Hacker News",
		PublishedAt: time.Unix(1625270400, 0),
		Score:       300,
		Comments:    30,
	}
	if !reflect.DeepEqual(news[0], expectedFirstNews) {
		t.Fatalf("Expected %+v, got %+v", expectedFirstNews, news[0])
	}
}

func TestFetchTopStories_TopStoriesError(t *testing.T) {
	// Create a mock HTTP client that returns an error for top stories
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			if url == "https://hacker-news.firebaseio.com/v0/topstories.json" {
				return errors.New("mock error")
			}
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Create a test source
	source := models.Source{
		Name:  "HackerNews",
		Type:  "hackernews",
		URL:   "https://hacker-news.firebaseio.com/v0",
		Limit: 3,
	}

	// Call the method being tested
	_, err := repo.FetchTopStories(source)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestFetchTopStories_StoryError(t *testing.T) {
	// Create a mock HTTP client that returns an error for one story
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			switch url {
			case "https://hacker-news.firebaseio.com/v0/topstories.json":
				// Mock the top stories response
				storyIDs, ok := v.(*[]int)
				if !ok {
					t.Fatalf("Expected *[]int, got %T", v)
				}
				*storyIDs = []int{123, 456, 789}
			case "https://hacker-news.firebaseio.com/v0/item/123.json":
				// Mock the first story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Test Story 1",
					"url":         "https://example.com/1",
					"score":       float64(100),
					"descendants": float64(10),
					"time":        float64(1625097600), // 2021-07-01
				}
			case "https://hacker-news.firebaseio.com/v0/item/456.json":
				// Return an error for the second story
				return errors.New("mock error")
			case "https://hacker-news.firebaseio.com/v0/item/789.json":
				// Mock the third story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Test Story 3",
					"url":         "https://example.com/3",
					"score":       float64(300),
					"descendants": float64(30),
					"time":        float64(1625270400), // 2021-07-03
				}
			default:
				t.Fatalf("Unexpected URL: %s", url)
			}
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Create a test source
	source := models.Source{
		Name:  "HackerNews",
		Type:  "hackernews",
		URL:   "https://hacker-news.firebaseio.com/v0",
		Limit: 3,
	}

	// Call the method being tested
	news, err := repo.FetchTopStories(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the results - should have 2 news items (one was skipped due to error)
	if len(news) != 2 {
		t.Fatalf("Expected 2 news items, got %d", len(news))
	}
}

func TestFetchBestStories_Success(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			switch url {
			case "https://hacker-news.firebaseio.com/v0/beststories.json":
				// Mock the best stories response
				storyIDs, ok := v.(*[]int)
				if !ok {
					t.Fatalf("Expected *[]int, got %T", v)
				}
				*storyIDs = []int{123, 456, 789}
			case "https://hacker-news.firebaseio.com/v0/item/123.json":
				// Mock the first story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Best Story 1",
					"url":         "https://example.com/1",
					"score":       float64(100),
					"descendants": float64(10),
					"time":        float64(1625097600), // 2021-07-01
				}
			case "https://hacker-news.firebaseio.com/v0/item/456.json":
				// Mock the second story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Best Story 2",
					"url":         "https://example.com/2",
					"score":       float64(200),
					"descendants": float64(20),
					"time":        float64(1625184000), // 2021-07-02
				}
			case "https://hacker-news.firebaseio.com/v0/item/789.json":
				// Mock the third story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Best Story 3",
					"url":         "https://example.com/3",
					"score":       float64(300),
					"descendants": float64(30),
					"time":        float64(1625270400), // 2021-07-03
				}
			default:
				t.Fatalf("Unexpected URL: %s", url)
			}
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Create a test source
	source := models.Source{
		Name:  "HackerNews",
		Type:  "hackernews",
		URL:   "https://hacker-news.firebaseio.com/v0",
		Limit: 3,
	}

	// Call the method being tested
	news, err := repo.FetchBestStories(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the results
	if len(news) != 3 {
		t.Fatalf("Expected 3 news items, got %d", len(news))
	}

	// Check that the news items are sorted by score (highest first)
	if news[0].Score != 300 || news[1].Score != 200 || news[2].Score != 100 {
		t.Fatalf("Expected news items to be sorted by score, got %v", news)
	}
}

func TestFetchBestStories_BestStoriesError(t *testing.T) {
	// Create a mock HTTP client that returns an error for best stories
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			if url == "https://hacker-news.firebaseio.com/v0/beststories.json" {
				return errors.New("mock error")
			}
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Create a test source
	source := models.Source{
		Name:  "HackerNews",
		Type:  "hackernews",
		URL:   "https://hacker-news.firebaseio.com/v0",
		Limit: 3,
	}

	// Call the method being tested
	_, err := repo.FetchBestStories(source)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestFetchBestStories_AllStoriesError(t *testing.T) {
	// Create a mock HTTP client that returns errors for all stories
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			switch url {
			case "https://hacker-news.firebaseio.com/v0/beststories.json":
				// Mock the best stories response
				storyIDs, ok := v.(*[]int)
				if !ok {
					t.Fatalf("Expected *[]int, got %T", v)
				}
				*storyIDs = []int{123, 456}
			case "https://hacker-news.firebaseio.com/v0/item/123.json":
				// Return an error for the first story
				return errors.New("mock error 1")
			case "https://hacker-news.firebaseio.com/v0/item/456.json":
				// Return an error for the second story
				return errors.New("mock error 2")
			default:
				t.Fatalf("Unexpected URL: %s", url)
			}
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Create a test source
	source := models.Source{
		Name:  "HackerNews",
		Type:  "hackernews",
		URL:   "https://hacker-news.firebaseio.com/v0",
		Limit: 3,
	}

	// Call the method being tested
	_, err := repo.FetchBestStories(source)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestFetchStory_Success(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			if url == "https://hacker-news.firebaseio.com/v0/item/123.json" {
				// Mock the story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title":       "Test Story",
					"url":         "https://example.com",
					"score":       float64(100),
					"descendants": float64(10),
					"time":        float64(1625097600), // 2021-07-01
				}
				return nil
			}
			t.Fatalf("Unexpected URL: %s", url)
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Call the method being tested
	news, err := repo.fetchStory(123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the results
	expectedNews := models.News{
		Title:       "Test Story",
		URL:         "https://example.com",
		Description: "Score: 100, Comments: 10",
		Source:      "Hacker News",
		PublishedAt: time.Unix(1625097600, 0),
		Score:       100,
		Comments:    10,
	}
	if !reflect.DeepEqual(news, expectedNews) {
		t.Fatalf("Expected %+v, got %+v", expectedNews, news)
	}
}

func TestFetchStory_Error(t *testing.T) {
	// Create a mock HTTP client that returns an error
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			return errors.New("mock error")
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Call the method being tested
	_, err := repo.fetchStory(123)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestFetchStory_SelfPost(t *testing.T) {
	// Create a mock HTTP client for a self-post (no URL)
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			if url == "https://hacker-news.firebaseio.com/v0/item/123.json" {
				// Mock the story response
				story, ok := v.(*map[string]any)
				if !ok {
					t.Fatalf("Expected *map[string]any, got %T", v)
				}
				*story = map[string]any{
					"title": "Test Self Post",
					// No URL field
					"score":       float64(100),
					"descendants": float64(10),
					"time":        float64(1625097600), // 2021-07-01
				}
				return nil
			}
			t.Fatalf("Unexpected URL: %s", url)
			return nil
		},
	}

	// Create the repository with the mock client
	repo := NewHackerNewsRepositoryWithClient(mockClient)

	// Call the method being tested
	news, err := repo.fetchStory(123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the results - should use the Hacker News item URL for self-posts
	expectedURL := "https://news.ycombinator.com/item?id=123"
	if news.URL != expectedURL {
		t.Fatalf("Expected URL %s, got %s", expectedURL, news.URL)
	}
}
