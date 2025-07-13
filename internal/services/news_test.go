package services

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
)

// MockHTTPClient is a mock implementation of the HTTP client for testing
type MockHTTPClient struct {
	GetJSONFunc        func(url string, v any) error
	GetFunc            func(url string) ([]byte, error)
	GetWithHeadersFunc func(url string, headers map[string]string) ([]byte, error)
}

// Get is a mock implementation of the Get method
func (m *MockHTTPClient) Get(url string) ([]byte, error) {
	if m.GetFunc != nil {
		return m.GetFunc(url)
	}
	return nil, errors.New("GetFunc not implemented")
}

// GetWithHeaders is a mock implementation of the GetWithHeaders method
func (m *MockHTTPClient) GetWithHeaders(url string, headers map[string]string) ([]byte, error) {
	if m.GetWithHeadersFunc != nil {
		return m.GetWithHeadersFunc(url, headers)
	}
	return nil, errors.New("GetWithHeadersFunc not implemented")
}

// GetJSON is a mock implementation of the GetJSON method
func (m *MockHTTPClient) GetJSON(url string, v any) error {
	if m.GetJSONFunc != nil {
		return m.GetJSONFunc(url, v)
	}
	return errors.New("GetJSONFunc not implemented")
}

// TestNewNewsService tests the NewNewsService function
func TestNewNewsService(t *testing.T) {
	// Create test sources
	sources := []models.Source{
		{
			Name:  "TestSource1",
			Type:  "hackernews",
			URL:   "https://example.com",
			Limit: 10,
		},
		{
			Name:      "TestSource2",
			Type:      "reddit",
			URL:       "https://example.org",
			Limit:     5,
			SubSource: "golang",
		},
	}

	// Create service
	service := NewNewsService(sources)

	// Check that the service was created correctly
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Check that the sources were stored correctly
	if len(service.sources) != len(sources) {
		t.Fatalf("Expected %d sources, got %d", len(sources), len(service.sources))
	}

	for i, source := range service.sources {
		if source.Name != sources[i].Name {
			t.Errorf("Expected source name %s, got %s", sources[i].Name, source.Name)
		}
	}

	// Check that the HTTP client was created
	if service.httpClient == nil {
		t.Fatal("Expected non-nil HTTP client")
	}
}

// TestFetchAllNews_Success tests the FetchAllNews function with successful fetches
func TestFetchAllNews_Success(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &MockHTTPClient{
		GetJSONFunc: func(url string, v any) error {
			if url == "https://hacker-news.firebaseio.com/v0/topstories.json" {
				// Mock the top stories response
				storyIDs, ok := v.(*[]int)
				if !ok {
					t.Fatalf("Expected *[]int, got %T", v)
				}
				*storyIDs = []int{123}
				return nil
			} else if url == "https://hacker-news.firebaseio.com/v0/item/123.json" {
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
			return errors.New("unexpected URL")
		},
		GetWithHeadersFunc: func(url string, headers map[string]string) ([]byte, error) {
			if url == "https://www.reddit.com/r/golang/hot.json?limit=1" {
				// Verify the User-Agent header is set correctly
				if userAgent, ok := headers["User-Agent"]; !ok || userAgent != "myredditbot/0.1 by u/ducminhgd" {
					return nil, errors.New("expected User-Agent header 'myredditbot/0.1 by u/ducminhgd'")
				}

				// Mock the Reddit response
				response := struct {
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
				}{
					Data: struct {
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
					}{
						Children: []struct {
							Data struct {
								Title       string  `json:"title"`
								URL         string  `json:"url"`
								Permalink   string  `json:"permalink"`
								Score       int     `json:"score"`
								NumComments int     `json:"num_comments"`
								Created     float64 `json:"created_utc"`
								Selftext    string  `json:"selftext"`
							} `json:"data"`
						}{
							{
								Data: struct {
									Title       string  `json:"title"`
									URL         string  `json:"url"`
									Permalink   string  `json:"permalink"`
									Score       int     `json:"score"`
									NumComments int     `json:"num_comments"`
									Created     float64 `json:"created_utc"`
									Selftext    string  `json:"selftext"`
								}{
									Title:       "Test Reddit Post",
									URL:         "https://example.com/reddit",
									Permalink:   "/r/golang/comments/123/test",
									Score:       200,
									NumComments: 20,
									Created:     1625184000, // 2021-07-02
									Selftext:    "Test post content",
								},
							},
						},
					},
				}
				return json.Marshal(response)
			}
			return nil, errors.New("unexpected URL")
		},
	}

	// Create test sources
	sources := []models.Source{
		{
			Name:  "HackerNews",
			Type:  "hackernews",
			URL:   "https://hacker-news.firebaseio.com/v0",
			Limit: 1,
		},
		{
			Name:      "RedditGo",
			Type:      "reddit",
			URL:       "https://www.reddit.com",
			Limit:     1,
			SubSource: "golang",
		},
	}

	// Create service with mock client
	service := &NewsService{
		httpClient: mockClient,
		sources:    sources,
	}

	// Call the method being tested
	result, err := service.FetchAllNews()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the results
	if len(result) != 2 {
		t.Fatalf("Expected 2 sources in result, got %d", len(result))
	}

	// Check HackerNews results
	hackerNews, ok := result["HackerNews"]
	if !ok {
		t.Fatal("Expected HackerNews in result")
	}
	if len(hackerNews) != 1 {
		t.Fatalf("Expected 1 HackerNews item, got %d", len(hackerNews))
	}
	expectedHackerNews := models.News{
		Title:       "Test Story",
		URL:         "https://example.com",
		Description: "Score: 100, Comments: 10",
		Source:      "Hacker News",
		PublishedAt: time.Unix(1625097600, 0),
		Score:       100,
		Comments:    10,
	}
	if !reflect.DeepEqual(hackerNews[0], expectedHackerNews) {
		t.Fatalf("Expected %+v, got %+v", expectedHackerNews, hackerNews[0])
	}

	// Check RedditGo results
	redditNews, ok := result["RedditGo"]
	if !ok {
		t.Fatal("Expected RedditGo in result")
	}
	if len(redditNews) != 1 {
		t.Fatalf("Expected 1 RedditGo item, got %d", len(redditNews))
	}
	expectedRedditNews := models.News{
		Title:       "Test Reddit Post",
		URL:         "https://example.com/reddit",
		Description: "Test post content",
		Source:      "Reddit",
		SubSource:   "golang",
		PublishedAt: time.Unix(1625184000, 0),
		Score:       200,
		Comments:    20,
	}
	if !reflect.DeepEqual(redditNews[0], expectedRedditNews) {
		t.Fatalf("Expected %+v, got %+v", expectedRedditNews, redditNews[0])
	}
}
