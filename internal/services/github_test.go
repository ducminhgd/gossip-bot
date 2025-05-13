package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ducminhgd/gossip-bot/internal/models"
	"github.com/google/go-github/v60/github"
)

// MockGithubClient is a mock implementation of the GitHub client for testing
type MockGithubClient struct {
	CreateIssueFunc func(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error)
}

// Issues is a mock implementation of the Issues service
type MockIssuesService struct {
	client *MockGithubClient
}

// Create is a mock implementation of the Create method
func (s *MockIssuesService) Create(ctx context.Context, owner, repo string, issue *github.IssueRequest) (*github.Issue, *github.Response, error) {
	return s.client.CreateIssueFunc(ctx, owner, repo, issue)
}

// TestNewGithubService tests the NewGithubService function
func TestNewGithubService(t *testing.T) {
	// Create a new service
	service := NewGithubService("test-token", "test-owner", "test-repo")

	// Check that the service was created correctly
	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	// Check that the owner and repo were stored correctly
	if service.owner != "test-owner" {
		t.Errorf("Expected owner to be 'test-owner', got '%s'", service.owner)
	}
	if service.repo != "test-repo" {
		t.Errorf("Expected repo to be 'test-repo', got '%s'", service.repo)
	}

	// Check that the client was created
	if service.client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// TestGenerateIssueContent tests the GenerateIssueContent function
func TestGenerateIssueContent(t *testing.T) {
	// Create a new service
	service := NewGithubService("test-token", "test-owner", "test-repo")

	// Create test news data
	newsMap := map[string][]models.News{
		"HackerNews": {
			{
				Title:       "Test HN Story 1",
				URL:         "https://example.com/hn1",
				Description: "Description 1",
				Source:      "Hacker News",
				PublishedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Score:       100,
				Comments:    10,
			},
			{
				Title:       "Test HN Story 2",
				URL:         "https://example.com/hn2",
				Description: "Description 2",
				Source:      "Hacker News",
				PublishedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Score:       200,
				Comments:    20,
			},
		},
		"RedditGo": {
			{
				Title:       "Test Reddit Post 1",
				URL:         "https://example.com/reddit1",
				Description: "Description 1",
				Source:      "Reddit",
				SubSource:   "golang",
				PublishedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Score:       300,
				Comments:    30,
			},
		},
		"EmptySource": {},
	}

	// Call the method being tested
	content, err := service.GenerateIssueContent(newsMap)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check the content
	// Get current date for comparison
	date := time.Now().UTC().Format("2006-01-02")

	// Expected content should have the date as a header
	if !strings.Contains(content, "# "+date) {
		t.Errorf("Expected content to contain date header '# %s', got: %s", date, content)
	}

	// Expected content should have HackerNews section
	if !strings.Contains(content, "## HackerNews") {
		t.Errorf("Expected content to contain HackerNews section, got: %s", content)
	}

	// Expected content should have RedditGo section
	if !strings.Contains(content, "## RedditGo") {
		t.Errorf("Expected content to contain RedditGo section, got: %s", content)
	}

	// Expected content should not have EmptySource section
	if strings.Contains(content, "## EmptySource") {
		t.Errorf("Expected content to not contain EmptySource section, got: %s", content)
	}

	// Expected content should have HackerNews stories
	if !strings.Contains(content, "1. [Test HN Story 1](https://example.com/hn1)") {
		t.Errorf("Expected content to contain HackerNews story 1, got: %s", content)
	}
	if !strings.Contains(content, "2. [Test HN Story 2](https://example.com/hn2)") {
		t.Errorf("Expected content to contain HackerNews story 2, got: %s", content)
	}

	// Expected content should have RedditGo post
	if !strings.Contains(content, "1. [Test Reddit Post 1](https://example.com/reddit1)") {
		t.Errorf("Expected content to contain RedditGo post 1, got: %s", content)
	}
}

// TestCreateIssue tests the CreateIssue function
func TestCreateIssue(t *testing.T) {
	// Skip this test since we can't easily mock the GitHub client
	t.Skip("Skipping TestCreateIssue because we can't easily mock the GitHub client")

	// The test would ideally verify:
	// 1. That the correct owner and repo are used
	// 2. That the issue title and body are correctly passed
	// 3. That the returned issue matches what we expect
}

// TestCreateIssue_Error tests the CreateIssue function with an error
func TestCreateIssue_Error(t *testing.T) {
	// Skip this test since we can't easily mock the GitHub client
	t.Skip("Skipping TestCreateIssue_Error because we can't easily mock the GitHub client")

	// The test would ideally verify:
	// 1. That when the GitHub API returns an error, our function also returns an error
	// 2. That the returned issue is nil
}
